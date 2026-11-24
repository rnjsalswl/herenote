// Package staticmap은 카카오 정적 지도 이미지를 서버에서 대신 받아 전달합니다.
//
// 카카오 REST 키는 서버 전용 키라서 브라우저에 노출되면 안 됩니다.
// 프론트는 /api/static-map 만 호출하고, 키는 이 서버의 환경변수에만 둡니다.
package staticmap

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const kakaoURL = "https://dapi.kakao.com/v2/maps/staticmap"

// 요청 크기 제한: 남이 이 엔드포인트로 큰 이미지를 마구 뽑아가지 못하게.
const (
	maxSide  = 800
	minSide  = 50
	maxLevel = 14
)

// Params는 검증을 통과한 요청 값입니다.
type Params struct {
	Lat, Lng float64
	Level    int
	W, H     int
}

// Parse는 쿼리스트링을 읽고 범위를 검사합니다.
func Parse(q url.Values) (Params, error) {
	var p Params
	var err error
	if p.Lat, err = strconv.ParseFloat(q.Get("lat"), 64); err != nil || p.Lat < -90 || p.Lat > 90 {
		return p, fmt.Errorf("lat이 올바르지 않습니다")
	}
	if p.Lng, err = strconv.ParseFloat(q.Get("lng"), 64); err != nil || p.Lng < -180 || p.Lng > 180 {
		return p, fmt.Errorf("lng가 올바르지 않습니다")
	}
	p.Level = intOr(q.Get("level"), 3)
	p.W = intOr(q.Get("w"), 400)
	p.H = intOr(q.Get("h"), 300)
	if p.Level < 1 || p.Level > maxLevel {
		return p, fmt.Errorf("level은 1~%d 사이여야 합니다", maxLevel)
	}
	if p.W < minSide || p.W > maxSide || p.H < minSide || p.H > maxSide {
		return p, fmt.Errorf("w, h는 %d~%d 사이여야 합니다", minSide, maxSide)
	}
	return p, nil
}

// UpstreamURL은 카카오로 보낼 주소를 만듭니다. 키는 여기서만 붙습니다.
func UpstreamURL(base, key string, p Params) string {
	center := fmt.Sprintf("%.6f,%.6f", p.Lng, p.Lat)
	v := url.Values{}
	v.Set("appkey", key)
	v.Set("center", center)
	v.Set("level", strconv.Itoa(p.Level))
	v.Set("size", fmt.Sprintf("%dx%d", p.W, p.H))
	v.Set("markers", "color:red|"+center)
	return base + "?" + v.Encode()
}

// Handler는 net/http 핸들러를 돌려줍니다. client가 nil이면 기본 클라이언트를 씁니다.
func Handler(key string, client *http.Client) http.HandlerFunc {
	return handler(kakaoURL, key, client)
}

func handler(base, key string, client *http.Client) http.HandlerFunc {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if key == "" {
			// 키가 없으면 프론트가 SVG 폴백을 그리도록 실패만 알려줌
			http.Error(w, "static map unavailable", http.StatusServiceUnavailable)
			return
		}
		p, err := Parse(r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := client.Get(UpstreamURL(base, key, p))
		if err != nil {
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// 카카오 쪽 에러 본문은 그대로 내보내지 않음 (키 관련 메시지가 섞일 수 있음)
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}

		ct := resp.Header.Get("Content-Type")
		if ct == "" {
			ct = "image/png"
		}
		w.Header().Set("Content-Type", ct)
		// 같은 좌표는 하루 동안 브라우저·CDN 캐시 → 카카오 호출 수 절약
		w.Header().Set("Cache-Control", "public, max-age=86400")
		io.Copy(w, io.LimitReader(resp.Body, 5<<20))
	}
}

func intOr(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1 // 범위 검사에서 걸러짐
	}
	return n
}
