package entry

import (
	"fmt"
	"log"

	"github.com/Xinrea/ffreplay/internal/component"
	"github.com/Xinrea/ffreplay/internal/data/fflogs"
	"github.com/Xinrea/ffreplay/internal/layer"
	"github.com/Xinrea/ffreplay/internal/model"
	"github.com/Xinrea/ffreplay/internal/tag"
	"github.com/Xinrea/ffreplay/pkg/object"
	"github.com/Xinrea/ffreplay/pkg/texture"
	"github.com/Xinrea/ffreplay/pkg/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/image/math/f64"
)

var Player = newArchetype(tag.GameObject, tag.Player, tag.PartyMember, tag.Buffable, component.Velocity, component.Sprite, component.Status)
var Enemy = newArchetype(tag.GameObject, tag.Enemy, tag.Buffable, component.Velocity, component.Sprite, component.Status)
var Background = newArchetype(tag.Background, component.Sprite)
var Camera = newArchetype(tag.Camera, component.Camera)
var Skill = newArchetype(tag.Skill, component.Skill)
var Timeline = newArchetype(tag.Timeline, component.Timeline)
var Marker = newArchetype(tag.Marker, component.Marker)
var Global = newArchetype(tag.Global, component.Global)

type archetype struct {
	components []donburi.IComponentType
}

func newArchetype(cs ...donburi.IComponentType) *archetype {
	return &archetype{
		components: cs,
	}
}

func (a *archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layer.Default,
		append(a.components, cs...)...,
	))
	return e
}

// boss gameID is unique in ffxiv, id is used in events
func NewEnemy(ecs *ecs.ECS, pos f64.Vec2, ringSize float64, gameID int64, id int64, name string, isBoss bool) *donburi.Entry {
	enemy := Enemy.Spawn(ecs)

	obj := object.NewPointObject(vector.NewVector(pos[0], pos[1]))
	textureRing := texture.NewTextureFromFile("asset/target_enemy.png")
	role := model.Boss
	if !isBoss {
		role = model.NPC
		textureRing = nil
	}
	component.Sprite.Set(enemy, &model.SpriteData{
		Texture:     textureRing,
		Scale:       ringSize,
		Face:        0,
		Object:      obj,
		Initialized: true,
	})
	component.Status.Set(enemy, &model.StatusData{
		GameID:   gameID,
		ID:       id,
		Name:     name,
		Role:     role,
		HP:       1,
		MaxHP:    1,
		Mana:     10000,
		MaxMana:  10000,
		BuffList: model.NewBuffList(),
		IsBoss:   isBoss,
	})

	return enemy
}

func NewPlayer(ecs *ecs.ECS, role model.RoleType, pos f64.Vec2, detail *fflogs.PlayerDetail) *donburi.Entry {
	player := Player.Spawn(ecs)
	var id int64 = 0
	name := "测试玩家"
	if detail != nil {
		id = detail.ID
		name = fmt.Sprintf("%s @%s", detail.Name, detail.Server)
		log.Println("Player:", name)
	}
	obj := object.NewPointObject(vector.NewVector(pos[0], pos[1]))
	// this scales target ring into size 50pixel, which means 1m in game
	component.Sprite.Set(player, &model.SpriteData{
		Texture:     texture.NewTextureFromFile("asset/target_normal.png"),
		Scale:       0.1842,
		Face:        0,
		Object:      obj,
		Initialized: true,
	})
	component.Status.Set(player, &model.StatusData{
		GameID:   -1,
		ID:       id,
		Name:     name,
		Role:     role,
		HP:       210000,
		MaxHP:    210000,
		Mana:     10000,
		MaxMana:  10000,
		BuffList: model.NewBuffList(),
	})

	return player
}

func NewMap(ecs *ecs.ECS, path string, offset f64.Vec2) *donburi.Entry {
	bg := Background.Spawn(ecs)
	obj := object.NewPointObject(vector.Vector(offset))
	component.Sprite.Set(bg, &model.SpriteData{
		Texture:     texture.NewTextureFromFile(path),
		Scale:       1,
		Face:        0,
		Object:      obj,
		Initialized: true,
	})

	return bg
}

func NewGlobal(ecs *ecs.ECS) *donburi.Entry {
	global := Global.Spawn(ecs)
	component.Global.Set(global, &model.GlobalData{
		Tick:  0,
		Speed: 10,
	})
	return global
}

func NewCamera(ecs *ecs.ECS) *donburi.Entry {
	camera := Camera.Spawn(ecs)
	component.Camera.Set(camera, &model.CameraData{
		ZoomFactor: 0,
		Rotation:   0,
	})
	return camera
}

func CastSkill(ecs *ecs.ECS, castTime int64, displayTime int64, gameSkill model.GameSkill) *donburi.Entry {
	newSkill := Skill.Spawn(ecs)
	component.Skill.Set(newSkill, &model.SkillData{
		Time: model.SkillTimeOption{
			StartTick:   GetTick(ecs),
			CastTime:    castTime,
			DisplayTime: displayTime,
		},
		GameSkill: gameSkill,
	})
	return newSkill
}

func NewTimeline(ecs *ecs.ECS, data *model.TimelineData) *donburi.Entry {
	timeline := Timeline.Spawn(ecs)
	component.Timeline.Set(timeline, data)
	return timeline
}

func NewMarker(ecs *ecs.ECS, markerType model.MarkerType, pos f64.Vec2) *donburi.Entry {
	marker := Marker.Spawn(ecs)
	component.Marker.Set(marker, &model.MarkerData{
		Type:     markerType,
		Position: pos,
	})
	return marker
}

func GetTick(ecs *ecs.ECS) int64 {
	return component.Global.Get(tag.Global.MustFirst(ecs.World)).Tick / 10
}

func GetSpeed(ecs *ecs.ECS) int64 {
	return component.Global.Get(tag.Global.MustFirst(ecs.World)).Speed
}
