package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapePage() {
	// ベースURLを設定
	baseURL := "某サイト" // 対象サイトのURL

	// HTTPリクエストを送信
	log.Printf("Fetching URL: %s", baseURL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyScraper/1.0)")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// ステータスコードの確認
	if res.StatusCode != http.StatusOK {
		log.Printf("HTTPエラー: %d %s", res.StatusCode, res.Status)
		return
	}

	// HTMLをパース
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// データを取得
	doc.Find(".col-sm-6.col-md-4.col-lg-4 .well").Each(func(i int, s *goquery.Selection) {
		fmt.Println("------------------------------")

		// 動画リンクを取得
		relativeLink, exists := s.Find("a.thumb-popu").Attr("href")
		fullLink := "動画リンクなし"
		if exists && relativeLink != "" {
			decodedLink, err := decodeURL(resolveURL(baseURL, relativeLink))
			if err == nil {
				fullLink = decodedLink
			}
		}
		fmt.Println("動画リンク:", fullLink)

		// サムネイルURLを取得
		imgSrc, exists := s.Find(".thumb-popu img.img-responsive").Attr("src")
		if !exists || imgSrc == "" {
			imgSrc = "サムネイルなし"
		}
		fmt.Println("サムネイル:", imgSrc)

		// 動画タイトルを取得
		videoTitle := s.Find(".video-title.title-truncate.m-t-5").Text()
		if videoTitle == "" {
			videoTitle = "タイトルなし"
		}
		fmt.Println("タイトル:", videoTitle)

		// 動画更新日を取得
		videoAdded := s.Find(".video-added").Text()
		if videoAdded == "" {
			videoAdded = "追加日なし"
		}
		fmt.Println("追加日:", strings.TrimSpace(videoAdded))

		// 視聴回数を取得
		videoViews := s.Find(".video-views.pull-left").Text()
		if videoViews == "" {
			videoViews = "再生回数なし"
		}
		fmt.Println("再生回数:", strings.TrimSpace(videoViews))

		// 評価を取得
		videoRating := s.Find(".video-rating.pull-right b").Text()
		if videoRating == "" {
			videoRating = "評価なし"
		}
		fmt.Println("評価:", strings.TrimSpace(videoRating))
	})
}

// resolveURL: ベースURLと相対URLを組み合わせてフルURLを作成
func resolveURL(base, relative string) string {
	u, err := url.Parse(base)
	if err != nil {
		return "URL解析エラー"
	}
	rel, err := url.Parse(relative)
	if err != nil {
		return "URL解析エラー"
	}
	return u.ResolveReference(rel).String()
}

// decodeURL: URLをデコード
func decodeURL(encoded string) (string, error) {
	decoded, err := url.QueryUnescape(encoded)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

func main() {
	ScrapePage()
}
