package stage07

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type Game struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type InputParams struct {
	Name string `json:"name" jsonschema:"description=the name of game"`
}

func GetGame(_ context.Context, params *InputParams) (string, error) {
	GameSet := []Game{
		{Name: "原神", Url: "https://ys.mihoyo.com/tool"},
		{Name: "王者荣耀", Url: "https://pvp.qq.com/"},
		{Name: "超级马里奥", Url: "https://www.nintendo.com/"},
	}
	for _, game := range GameSet {
		if game.Name == params.Name {
			return game.Url, nil
		}
	}
	return "Unknown game", nil
}

func CreateTool() tool.InvokableTool {
	getGameTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_game",
		Desc: "get a game url by name",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": {
					Type:     schema.String,
					Desc:     "game's name",
					Required: true,
				},
			},
		),
	}, GetGame)
	return getGameTool
}
