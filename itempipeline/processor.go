package itempipeline

import "go-crawler/base"

type ProcessItem func(item base.Item) (result base.Item, err error)
