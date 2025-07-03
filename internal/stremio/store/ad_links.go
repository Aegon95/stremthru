package stremio_store

import (
	"net/http"
	"time"

	"github.com/MunifTanjim/stremthru/core"
	"github.com/MunifTanjim/stremthru/internal/cache"
	"github.com/MunifTanjim/stremthru/internal/config"
	"github.com/MunifTanjim/stremthru/internal/context"
	"github.com/MunifTanjim/stremthru/internal/shared"
	stremio_store_webdl "github.com/MunifTanjim/stremthru/internal/stremio/store/webdl"
	"github.com/MunifTanjim/stremthru/stremio"
)

var adLinksCache = cache.NewCache[[]stremio.MetaVideo](&cache.CacheConfig{
	Lifetime: 10 * time.Minute,
	Name:     "stremio:store:ad:links",
})

func getADLinksCacheKey(idPrefix, storeToken string) string {
	return idPrefix + storeToken
}

func getADWebDLsMeta(r *http.Request, ctx *context.StoreContext, idr *ParsedId) stremio.Meta {
	released := time.Now().UTC()

	meta := stremio.Meta{
		Id:          getWebDLsMetaId(idr.getStoreCode()),
		Type:        ContentTypeOther,
		Name:        "Web Downloads",
		Description: "Web Downloads from AllDebrid",
		Released:    &released,
		Videos:      []stremio.MetaVideo{},
	}
	cacheKey := getADLinksCacheKey(getIdPrefix(idr.getStoreCode()), ctx.StoreAuthToken)
	if !adLinksCache.Get(cacheKey, &meta.Videos) {
		params := &stremio_store_webdl.ListWebDLsParams{}
		params.APIKey = ctx.StoreAuthToken
		res, err := stremio_store_webdl.ListWebDLs(params, idr.storeName)
		if err != nil {
			log.Error("failed to list webdls", "error", err, "store", idr.storeCode)
			return meta
		}

		storeName := ctx.Store.GetName()
		shouldCreateProxyLink := config.StoreContentProxy.IsEnabled(string(storeName)) && ctx.StoreAuthToken == config.StoreAuthToken.GetToken(ctx.ProxyAuthUser, string(storeName)) && ctx.IsProxyAuthorized
		tunnelType := config.StoreTunnel.GetTypeForStream(string(ctx.Store.GetName()))
		idPrefix := getWebDLsMetaIdPrefix(idr.getStoreCode())

		for i := range res.Items {
			dl := &res.Items[i]
			if len(dl.Files) == 0 {
				continue
			}
			file := dl.Files[0]
			if !core.HasVideoExtension(file.Name) {
				continue
			}
			stream := stremio.Stream{
				URL: file.Link,
				BehaviorHints: &stremio.StreamBehaviorHints{
					VideoSize: file.Size,
					Filename:  file.Name,
				},
			}
			videoTitle := getMetaPreviewDescriptionForWebDL("", dl.Name, true) + "\n📄 " + file.Name
			if shouldCreateProxyLink {
				if proxyLink, err := shared.CreateProxyLink(r, stream.URL, nil, tunnelType, 12*time.Hour, ctx.ProxyAuthUser, ctx.ProxyAuthPassword, true, file.Name); err == nil {
					stream.URL = proxyLink
					videoTitle = "✨ " + videoTitle
				} else {
					log.Error("failed to create proxy link, skipping file", "error", err, "store", storeName.Code(), "filename", file.Name)
					continue
				}
			}
			video := stremio.MetaVideo{
				Id:       idPrefix + dl.Id,
				Title:    videoTitle,
				Released: dl.AddedAt,
				Streams:  []stremio.Stream{stream},
				Episode:  -1,
				Season:   -1,
			}
			meta.Videos = append(meta.Videos, video)
		}

		rdDownloadsCache.Add(cacheKey, meta.Videos)
	}
	return meta
}
