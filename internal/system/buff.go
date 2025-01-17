package system

import (
	"github.com/Xinrea/ffreplay/internal/component"
	"github.com/Xinrea/ffreplay/internal/entry"
	"github.com/Xinrea/ffreplay/internal/model"
	"github.com/Xinrea/ffreplay/internal/tag"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type RemoveCallbackDB map[int64]func(*model.Buff, *ecs.ECS, *donburi.Entry, *donburi.Entry)

var buffRemoveDB RemoveCallbackDB = make(RemoveCallbackDB)

func (s *System) BuffUpdate(ecs *ecs.ECS, tick int64) {
	global := entry.GetGlobal(ecs)
	if !global.Loaded.Load() {
		return
	}

	for e := range tag.Buffable.Iter(ecs.World) {
		component.Status.Get(e).BuffList.Update(tick)
	}

	if entry.GetGlobal(s.ecs).ReplayMode {
		return
	}

	for e := range tag.Buffable.Iter(ecs.World) {
		component.Status.Get(e).BuffList.UpdateExpire(tick)
	}
}
