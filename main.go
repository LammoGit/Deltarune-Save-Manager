package main

import (
	"context"
	"runtime"
	"path"
	"fmt"
	"os"
	"slices"
	"maps"

	"github.com/urfave/cli/v3"

	"dsm/saves"
)

func getSavesFolderPath() (dirPath string, err error) {
	if runtime.GOOS == "windows" {
		dirPath, err = os.UserCacheDir()
		dirPath = path.Join(dirPath, "DELTARUNE")
	}
	return
}

func dirExist(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

var (
	dirPath string
)

func main() {
	dirPath, err := getSavesFolderPath()
	if err != nil {
		fmt.Println(err)
		return
	}

	manager, err := saves.NewSaveManager(path.Join(dirPath, "save-manager"), dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name: "save",
				Aliases: []string{},
				Usage: "commands for managing save files",
				Commands: []*cli.Command{
					{
						Name: "create",
						Aliases: []string{},
						Usage: "create new save file with standard properties",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Create(
								cmd.StringArg("save name"),
								cmd.IntArg("chapter"),
							)
							if err != nil {
								fmt.Println(err)
							}
							return nil
						},
					},
					{
						Name: "remove",
						Aliases: []string{},
						Usage: "remove a save file with or without connected slots",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "remove-slots",
								Aliases: []string{"slots", "cascade"},
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Remove(
								cmd.StringArg("save name"),
								cmd.IntArg("chapter"),
								cmd.Bool("remove-slots"),
							)
							if err != nil {
								fmt.Println(err)
							}
							return nil	
						},
					},
					{
						Name: "set",
						Aliases: []string{},
						Usage: "set a save file in given save slot",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "erase-unmanaged",
								Aliases: []string{"unmanaged"},
							},
						},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save name",
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
								cmd.StringArg("save name"),
								cmd.IntArg("chapter"),
								cmd.IntArg("slot"),
								cmd.Bool("erase-unmanaged"),
							)
							if err != nil {
								fmt.Println(err)
							}
							return nil
						},
					},
					{
						Name: "rename",
						Aliases: []string{},
						Usage: "rename a save file",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save name from",
							},
							&cli.StringArg{
								Name: "save name to",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Rename(
								cmd.StringArg("save name from"),
								cmd.StringArg("save name to"),
								cmd.IntArg("chapter"),
							)
							if err != nil {
								fmt.Println(err)
							}
							return nil	
						},
					},
					{
						Name: "swap",
						Aliases: []string{},
						Usage: "swap names of two save files",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "first save name",
							},
							&cli.StringArg{
								Name: "second save name",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Swap(
								cmd.StringArg("first save name"),
								cmd.StringArg("second save name"),
								cmd.IntArg("chapter"),
							)
							if err != nil {
								fmt.Println(err)
							}
							return nil
						},
					},
					{
						Name: "copy",
						Aliases: []string{},
						Usage: "create a copy of a save file",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save name from",
							},
							&cli.StringArg{
								Name: "save name to",
							},
							&cli.IntArg{
								Name: "chapter",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := manager.Copy(
								cmd.StringArg("save name from"),
								cmd.StringArg("save name to"),
								cmd.IntArg("chapter"),
							)
							if err != nil {
								fmt.Println(err)
							}
							return nil	
						},
					},
					{
						Name: "edit",
						Aliases: []string{},
						Usage: "change properties of a save file",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return nil	
						},
					},
				},
			},
			{
				Name: "slot",
				Aliases: []string{},
				Usage: "commands for managing save slots",
				Commands: []*cli.Command{
					{
						Name: "save",
						Aliases: []string{},
						Usage: "create a save file of a slot",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "save name",
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
								cmd.StringArg("save name"),
								cmd.IntArg("chapter"),
								cmd.IntArg("slot"),
							)
							if err != nil {
								fmt.Println(err)
							}
							return nil
						},
					},
					{
						Name: "unset",
						Aliases: []string{},
						Usage: "remove save slot, but keep save file",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name: "erase-unmanaged",
								Usage: "remove slots without a linked save file",
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
							if err != nil {
								fmt.Println(err)
							}
							return nil
						},
					},
				},
			},
			{
				Name: "saves",
				Aliases: []string{},
				Usage: "list all save files",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					for id := range manager.Saves {
						fmt.Printf("Chapter %d - %s\n", id.Chapter, id.Name)
					}
					return nil	
				},
			},
			{
				Name: "slots",
				Aliases: []string{},
				Usage: "list all save slots",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					keys := slices.Collect(maps.Keys(manager.Slots))
					slices.SortFunc(
						keys,
						func(a, b saves.SlotID) int {
							if a.Chapter == b.Chapter && a.Slot == b.Slot { return 0}
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
						case saves.Save1:
							playerName = s.PlayerName
							charName = s.CharName
						case saves.Save2:
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

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}
}
