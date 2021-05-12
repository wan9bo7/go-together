package task

import (
	"github.com/gocolly/colly/v2"
	"together/blog_server/pkg/assets"
	"together/global"
	pb "together/proto"
)

type BlogJob struct {
	Domains []string
}

func (b *BlogJob) SetDomain(domains []string) {
	b.Domains = domains
}

func (b *BlogJob) Run() {
	for _, domain := range b.Domains {
		var menus []*pb.GetListReply_Data
		switch domain {
		case global.BlogServer.IxugoDomain:
			menus = getIxugo(domain)
		case global.BlogServer.WangboDomain:
			menus = getWangbo(domain)
		}
		if menus != nil {

		}
	}
}

func getIxugo(url string) []*pb.GetListReply_Data {
	// TODO 识别链接中的域名作为参数填入下方
	a := assets.GetInstance()
	menus := make([]*pb.GetListReply_Data, 0, 10)
	a.OnHTML("main section", func(e *colly.HTMLElement) {
		e.ForEach("article", func(i int, h *colly.HTMLElement) {
			const prefix = "header div"
			art := pb.GetListReply_Data{
				Img:         "",
				Title:       h.ChildText(prefix + " h2 a"),
				Description: "",
				CreateAt:    h.ChildText(prefix + " footer time"),
				Tags:        h.ChildTexts(prefix + " header a"),
				Category:    "",
				Link:        url + h.ChildAttr(prefix+" h2 a", "href"),
			}
			menus = append(menus, &art)
		})
	})
	a.Visit(url)
	return menus
}

func getWangbo(url string) []*pb.GetListReply_Data {
	// TODO 识别链接中的域名作为参数填入下方
	a := assets.GetInstance()
	menus := make([]*pb.GetListReply_Data, 0, 10)
	a.OnHTML(".recent-posts", func(e *colly.HTMLElement) {
		e.ForEach(".recent-post-item", func(i int, h *colly.HTMLElement) {
			const website = "https://chenyunxin.cn"
			art := pb.GetListReply_Data{
				Img:         h.ChildAttr(".post_cover a img", "data-original"),
				Title:       h.ChildAttr(".post_cover a", "title"),
				Description: "",
				CreateAt:    h.ChildText(".recent-post-info div time"),
				Tags:        []string{},
				Category:    h.ChildText(".article-meta__categories"),
				Link:        website + h.ChildAttr(".post_cover a", "href"),
			}
			menus = append(menus, &art)
		})
	})
	a.Visit(url)
	return menus
}
