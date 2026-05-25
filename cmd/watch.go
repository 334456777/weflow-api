package cmd

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var watchReconnect bool

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "订阅 SSE 实时消息（§2 push）",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		r := newRenderer()
		ctx := cmd.Context()

		q := url.Values{}
		if c.Token() != "" {
			q.Set("access_token", c.Token())
		}

		var attempt int
		for {
			events, errs := c.StreamEvents(ctx, "/api/v1/push/messages", q)

			for e := range events {
				if err := r.Event(e); err != nil {
					return err
				}
				attempt = 0
			}

			var streamErr error
			select {
			case streamErr = <-errs:
			default:
			}

			if ctx.Err() != nil {
				return nil
			}
			if !watchReconnect {
				return streamErr
			}

			attempt++
			wait := time.Duration(math.Min(30, math.Pow(2, float64(attempt)))) * time.Second
			if streamErr != nil {
				fmt.Fprintf(os.Stderr, "连接中断: %v；%s 后重连\n", streamErr, wait)
			} else {
				fmt.Fprintf(os.Stderr, "连接关闭；%s 后重连\n", wait)
			}
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return nil
			}
		}
	},
}

func init() {
	watchCmd.Flags().BoolVar(&watchReconnect, "reconnect", false, "断线自动重连（指数退避，最长 30s）")
	rootCmd.AddCommand(watchCmd)
}
