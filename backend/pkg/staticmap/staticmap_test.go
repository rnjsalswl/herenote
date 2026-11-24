package staticmap

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestParseRejectsBadInput(t *testing.T) {
	cases := []string{
		"lng=127",                     // lat 없음
		"lat=91&lng=127",              // 위도 범위 밖
		"lat=37&lng=200",              // 경도 범위 밖
		"lat=37&lng=127&level=0",      // 줌 범위 밖
		"lat=37&lng=127&w=5000&h=300", // 너무 큰 이미지
		"lat=37&lng=127&w=abc",        // 숫자 아님
	}
	for _, c := range cases {
		q, _ := url.ParseQuery(c)
		if _, err := Parse(q); err == nil {
			t.Errorf("%q: 에러가 나야 하는데 통과함", c)
		}
	}
}

func TestParseDefaults(t *testing.T) {
	q, _ := url.ParseQuery("lat=37.58&lng=126.97")
	p, err := Parse(q)
	if err != nil {
		t.Fatal(err)
	}
	if p.Level != 3 || p.W != 400 || p.H != 300 {
		t.Errorf("기본값이 다름: %+v", p)
	}
}

func TestHandlerKeepsKeyServerSide(t *testing.T) {
	var gotKey string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.URL.Query().Get("appkey")
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte("PNG"))
	}))
	defer upstream.Close()

	h := handler(upstream.URL, "secret-rest-key", upstream.Client())
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest("GET", "/api/static-map?lat=37.58&lng=126.97&level=3&w=400&h=300", nil))

	if rec.Code != 200 || rec.Body.String() != "PNG" {
		t.Fatalf("응답이 이상함: %d %q", rec.Code, rec.Body.String())
	}
	if gotKey != "secret-rest-key" {
		t.Errorf("카카오로 키가 안 넘어감: %q", gotKey)
	}
	if strings.Contains(rec.Body.String()+rec.Header().Get("Location"), "secret-rest-key") {
		t.Error("응답에 키가 섞여 나감")
	}
	if rec.Header().Get("Cache-Control") == "" {
		t.Error("캐시 헤더 없음")
	}
}

func TestHandlerWithoutKey(t *testing.T) {
	rec := httptest.NewRecorder()
	Handler("", nil)(rec, httptest.NewRequest("GET", "/api/static-map?lat=37&lng=127", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("키 없으면 503이어야 함: %d", rec.Code)
	}
}

func TestHandlerHidesUpstreamError(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid appkey secret-rest-key", http.StatusUnauthorized)
	}))
	defer upstream.Close()

	rec := httptest.NewRecorder()
	handler(upstream.URL, "secret-rest-key", upstream.Client())(rec, httptest.NewRequest("GET", "/api/static-map?lat=37&lng=127", nil))
	if rec.Code != http.StatusBadGateway || strings.Contains(rec.Body.String(), "secret") {
		t.Errorf("카카오 에러가 그대로 노출됨: %d %q", rec.Code, rec.Body.String())
	}
}
