package client

import (
	"context"
	"fmt"
	"github.com/cyberax/go-dd-service-base/utils"
	"github.com/cyberax/go-dd-service-base/visibility"
	"github.com/spf13/cobra"
	"github.com/twitchtv/twirp"
	"net/http"
	"xyzzy/gen/xyzzy"
)

func makeClient(cmd *cobra.Command) (xyzzy.XyzzyData, error) {
	serverUrl := utils.GetFlagS(cmd, "server")

	token := utils.GetFlagS(cmd, "token")
	//if token == "" {
	//	return nil, fmt.Errorf("you need to specify the token or XYZZY_TOKEN env var")
	//}

	hooks := &twirp.ClientHooks{}
	hooks.RequestPrepared = func(ctx context.Context, request *http.Request) (context.Context, error) {
		if token != "" {
			request.Header.Add("Authorization", "Bearer "+token)
		}
		return ctx, nil
	}

	cli := visibility.WrapTwirpClientDef(http.DefaultClient, "xyzzy_cli")
	client := xyzzy.NewXyzzyDataProtobufClient(serverUrl, cli, twirp.WithClientHooks(hooks))
	return client, nil
}

func MakePingCmd() *cobra.Command {
	var cmdPing = &cobra.Command{
		Use:           "ping",
		Short:         "Run Ping",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := makeClient(cmd)
			if err != nil {
				return err
			}
			return doPing(context.TODO(), cli)
		},
	}

	return cmdPing
}

func doPing(ctx context.Context, cli xyzzy.XyzzyData) error {
	_, err := cli.Ping(ctx, &xyzzy.PingRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("OK\n")
	return nil
}
