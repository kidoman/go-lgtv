package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kidoman/go-lgtv/control"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

func main() {
	clientKeyFlag := cli.StringFlag{
		Name:     "client-key",
		Required: true,
	}

	app := cli.App{
		Name:  "lgtv",
		Usage: "LG TV control",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "ip",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "mask",
				Value: "255.255.255.0",
			},
			&cli.StringFlag{
				Name:     "mac",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "initialize the connection",
				Action: initAction,
			},
			{
				Name:   "on",
				Usage:  "switch tv on",
				Action: onAction,
			},
			{
				Name:  "off",
				Usage: "switch tv off",
				Flags: []cli.Flag{
					&clientKeyFlag,
				},
				Action: offAction,
			},
			{
				Name:  "input-connected",
				Usage: "whether input is connected",
				Flags: []cli.Flag{
					&clientKeyFlag,
					&cli.StringFlag{
						Name: "input-id",
					},
				},
				Action: inputConnectionAction,
			},
			{
				Name:  "input-switch",
				Usage: "switch to input",
				Flags: []cli.Flag{
					&clientKeyFlag,
					&cli.StringFlag{
						Name: "input-id",
					},
				},
				Action: inputSwitchAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type connectOptions struct {
	ip        string
	mask      string
	mac       string
	clientKey string
}

func parseFlags(cCtx *cli.Context) *connectOptions {
	return &connectOptions{
		ip:        cCtx.String("ip"),
		mask:      cCtx.String("mask"),
		mac:       cCtx.String("mac"),
		clientKey: cCtx.String("client-key"),
	}
}

func initAction(cCtx *cli.Context) error {
	ip := cCtx.String("ip")

	tv, err := control.NewTV(ip, "", "")
	if err != nil {
		return err
	}

	key, err := tv.Connect("", 1000)
	if err != nil {
		return err
	}

	log.Printf("client key is %q", key)

	return nil
}

func onAction(cCtx *cli.Context) error {
	o := parseFlags(cCtx)

	tv, err := control.NewTV(o.ip, o.mac, o.mask)
	if err != nil {
		return err
	}

	return tv.TurnOn()
}

func offAction(cCtx *cli.Context) error {
	o := parseFlags(cCtx)

	tv, err := control.NewTV(o.ip, "", "")
	if err != nil {
		return err
	}

	_, err = tv.Connect(o.clientKey, 1000)
	if err != nil {
		return err
	}
	defer tv.Disconnect() //nolint:errcheck

	return tv.TurnOff()
}

func inputConnectionAction(cCtx *cli.Context) error {
	o := parseFlags(cCtx)

	tv, err := control.NewTV(o.ip, "", "")
	if err != nil {
		return err
	}

	_, err = tv.Connect(o.clientKey, 1000)
	if err != nil {
		return err
	}
	defer tv.Disconnect() //nolint:errcheck

	inputs, err := tv.ListExternalInputs()
	if err != nil {
		return err
	}

	inputID := cCtx.String("input-id")
	if idx := slices.IndexFunc(inputs, func(i control.Input) bool {
		return i.ID == inputID
	}); idx == -1 {
		fmt.Println("not-connected")
	} else {
		fmt.Println("connected")
	}

	return nil
}

func inputSwitchAction(cCtx *cli.Context) error {
	o := parseFlags(cCtx)

	tv, err := control.NewTV(o.ip, "", "")
	if err != nil {
		return err
	}

	_, err = tv.Connect(o.clientKey, 1000)
	if err != nil {
		return err
	}
	defer tv.Disconnect() //nolint:errcheck

	return tv.SwitchInput(cCtx.String("input-id"))
}
