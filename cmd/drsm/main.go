package main

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"

	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/LammoGit/Deltarune-Save-Manager/saves"
	"github.com/LammoGit/Deltarune-Save-Manager/utils"
)

func main() {
	// Get Deltarune saves folder
	dirPath, err := utils.GetSavesFolderPath()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a save manager object
	manager, err := saves.NewSaveManager(filepath.Join(dirPath, "save-manager"), dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// CLI-utility description
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "save",
				Aliases: []string{},
				Usage:   "commands for managing save files",
				Commands: []*cli.Command{
					{
						Name:  "create",
						Usage: "create new save file with standard properties",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save_name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Create(
								cmd.StringArg("save_name"),
								cmd.IntArg("chapter"),
								cmd.Bool("sideb"),
							)
							return err
						},
					},
					{
						Name:  "info",
						Usage: "show save's contents",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save_name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							saveID := saves.SaveID{
								Name:    cmd.StringArg("save_name"),
								Chapter: cmd.IntArg("chapter"),
								SideB:   cmd.Bool("sideb"),
							}
							save, ok := manager.Saves[saveID]
							if !ok {
								return nil
							}
							fmt.Println(save)
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove a save file with or without connected slots",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "remove-slots",
								Aliases: []string{"slots", "cascade"},
							},
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save_name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Remove(
								cmd.StringArg("save_name"),
								cmd.IntArg("chapter"),
								cmd.Bool("sideb"),
								cmd.Bool("remove-slots"),
							)
							return err
						},
					},
					{
						Name:  "rename",
						Usage: "rename a save file",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save_name from",
							},
							&cli.StringArg{
								Name: "save_name to",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Rename(
								cmd.StringArg("save_name from"),
								cmd.StringArg("save_name to"),
								cmd.IntArg("chapter"),
								cmd.Bool("sideb"),
							)
							return err
						},
					},
					{
						Name:  "swap",
						Usage: "swap names of two save files",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "first save_name",
							},
							&cli.StringArg{
								Name: "second save_name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Swap(
								cmd.StringArg("first save_name"),
								cmd.StringArg("second save_name"),
								cmd.IntArg("chapter"),
								cmd.Bool("sideb"),
							)
							return err
						},
					},
					{
						Name:  "copy",
						Usage: "create a copy of a save file",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save_name from",
							},
							&cli.StringArg{
								Name: "save_name to",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Copy(
								cmd.StringArg("save_name from"),
								cmd.StringArg("save_name to"),
								cmd.IntArg("chapter"),
								cmd.Bool("sideb"),
							)
							return err
						},
					},
					{
						Name:  "edit",
						Usage: "change properties of a save file",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return nil
						},
					},
				},
			},
			{
				Name:    "slot",
				Aliases: []string{},
				Usage:   "commands for managing save slots",
				Commands: []*cli.Command{
					{
						Name:  "save",
						Usage: "create a save file of a slot",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save_name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
							&cli.IntArg{
								Name: "slot",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.SaveSlot(
								cmd.StringArg("save_name"),
								cmd.IntArg("chapter"),
								cmd.IntArg("slot"),
								cmd.Bool("sideb"),
							)
							return err
						},
					},
					{
						Name:  "info",
						Usage: "show slot's contents",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.IntArg{
								Name: "chapter",
							},
							&cli.IntArg{
								Name: "slot",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							slotID := saves.SlotID{
								Chapter: cmd.IntArg("chapter"),
								Slot:    cmd.IntArg("slot"),
								SideB:   cmd.Bool("sideb"),
							}
							save, ok := manager.Slots[slotID]
							if !ok {
								return nil
							}
							fmt.Println(save)
							return nil
						},
					},
					{
						Name:  "set",
						Usage: "set a save file in given save slot",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "erase-unmanaged",
								Aliases: []string{"unmanaged"},
							},
							&cli.BoolFlag{
								Name: "sideb",
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save_name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
							&cli.IntArg{
								Name: "slot",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.SetSlot(
								cmd.StringArg("save_name"),
								cmd.IntArg("chapter"),
								cmd.IntArg("slot"),
								cmd.Bool("sideb"),
								cmd.Bool("erase-unmanaged"),
							)
							return err
						},
					},
					{
						Name:  "unset",
						Usage: "remove save slot, but keep save file",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "erase-unmanaged",
								Usage:   "remove slots without a linked save file",
								Aliases: []string{"unmanaged"},
							},
						},
						Arguments: []cli.Argument{
							&cli.IntArg{
								Name: "chapter",
							},
							&cli.IntArg{
								Name: "slot",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.UnsetSlot(
								cmd.IntArg("chapter"),
								cmd.IntArg("slot"),
								cmd.Bool("erase-unmanaged"),
							)
							return err
						},
					},
				},
			},
			{
				Name:  "saves",
				Usage: "list all save files",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					for id := range manager.Saves {
						fmt.Printf("Chapter %d - %s\n", id.Chapter, id.Name)
					}
					return nil
				},
			},
			{
				Name:  "slots",
				Usage: "list all save slots",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					keys := slices.Collect(maps.Keys(manager.Slots))
					slices.SortFunc(
						keys,
						func(a, b saves.SlotID) int {
							if a.Chapter == b.Chapter && a.Slot == b.Slot {
								return 0
							}
							if a.Chapter < b.Chapter ||
								(a.Chapter == b.Chapter && a.Slot < b.Slot) {
								return -1
							}
							return 1
						},
					)
					for _, id := range keys {
						slot := manager.Slots[id]

						var playerName string
						var charName string
						switch s := slot.(type) {
						case *saves.Save1:
							playerName = s.PlayerName
							charName = s.CharName
						case *saves.Save2:
							playerName = s.PlayerName
							charName = s.CharName
						}
						fmt.Printf(
							"Chapter:%d Slot:%d Player:%s Character:%s\n",
							id.Chapter,
							id.Slot,
							playerName,
							charName,
						)
					}
					return nil
				},
			},
		},
	}

	// Run CLI-utility
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}
}
