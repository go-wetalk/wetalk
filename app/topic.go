package app

import (
	"appsrv/pkg/db"
	"appsrv/schema"
	"appsrv/service"
	"net/http"

	"github.com/kataras/muxie"
	"github.com/spf13/cast"
)

// Topic 话题
type Topic struct{}

// List 取出话题列表
func (Topic) List(w http.ResponseWriter, r *http.Request) {
	input := schema.TopicListInput{}
	input.Size = 20
	if p := r.URL.Query().Get("p"); p != "" {
		input.Page = cast.ToUint(p)
	}

	ts, _ := service.Topic{}.ListWithRankByScore(db.DB, input)
	muxie.Dispatch(w, muxie.JSON, ts)
}
