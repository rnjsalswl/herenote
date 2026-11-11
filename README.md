# herenote

**그 장소에 직접 가야만 읽고 쓸 수 있는 위치 기반 방명록 PWA**

산 정상, 공원, 카페 같은 장소마다 방명록이 있고, GPS로 현장에 있다는 게 확인돼야만 그곳에 남겨진 글을 열람하고 새 글을 남길 수 있습니다.

<!-- TODO: 스크린샷 또는 GIF (홈 / 장소 상세 / 위치 인증 / 방명록 작성) -->

## 주요 기능

- **위치 인증**: 현재 좌표를 서버로 보내면 PostGIS로 장소 반경(50m~500m) 안에 있는지 검증하고, 24시간 유효한 장소 전용 인증 토큰을 발급
- **방명록**: 인증 토큰이 있어야만 해당 장소의 글을 열람하고 작성 가능 (최대 500자)
- **회원 / 로그인**: 닉네임과 비밀번호로 가입하고, bcrypt 해시 저장, JWT(HS256, 30일) 발급
- **프로필**: 내가 다녀간 장소, 작성한 글 수, 뱃지(BLUE / GOLD / SPECIAL)
- **장소 탐색**: 카카오맵 SDK로 지도, 주변 장소, 정적 지도 썸네일 표시
- **PWA**: 홈 화면 설치, 서비스 워커 캐싱
- **어드민**: 관리자 계정만 장소 등록, 뱃지 부여 가능

## 기술 스택

| 영역 | 사용 기술 |
|---|---|
| Backend | Go, Gin, pgx (PostgreSQL), PostGIS, golang-jwt, bcrypt |
| Frontend | Vue 3 (Composition API), Vue Router, Vite, vite-plugin-pwa |
| 외부 API | Kakao Maps JS SDK, Kakao Static Map |
| 배포 | Railway |

## 아키텍처

```
[Vue PWA] ──HTTPS──▶ [Gin Router]
                        │  middleware: CORS / UserAuth(JWT) / LocationAuth / AdminOnly
                        ▼
                     handler  ──▶  service  ──▶  repository  ──▶  PostgreSQL + PostGIS
```

```
backend/
├── main.go                 # 의존성 조립, 라우팅
├── config/                 # 환경변수 로드 (필수 시크릿이 없으면 서버가 기동하지 않음)
├── internal/
│   ├── auth/               # JWT 발급·검증
│   ├── middleware/         # 인증·인가 미들웨어
│   ├── handler/            # HTTP 요청/응답
│   ├── service/            # 비즈니스 로직
│   ├── repository/         # SQL (pgx)
│   └── model/
└── pkg/gps/                # Haversine 거리 계산
frontend/src/
├── views/  components/  composables/  stores/  router/
```

## 인증 흐름

인증이 두 겹으로 되어 있습니다. **"누구인가"**(JWT)와 **"지금 거기 있는가"**(위치 토큰)를 분리했습니다.

1. `POST /api/auth/login` → JWT 발급
2. 장소에서 `POST /api/places/:id/verify` (JWT + 현재 좌표) → 서버가 PostGIS로 반경 검증 → 성공하면 위치 토큰(UUID, 24h)을 `verified_sessions`에 저장하고 응답
3. 방명록 API는 `X-Location-Token` 헤더로 장소별 인증 세션을 확인

위치 판정을 클라이언트가 아니라 **서버에서** 하기 때문에, 프론트 코드를 조작해서 "인증됨" 상태로 바꾸는 우회가 통하지 않습니다.

## API

| Method | Path | 인증 | 설명 |
|---|---|---|---|
| POST | `/api/users` | - | 회원가입 |
| POST | `/api/auth/login` | - | 로그인 |
| GET | `/api/users/:id` | - | 유저 조회 |
| GET | `/api/places` | - | 장소 목록 |
| GET | `/api/places/nearby?lat=&lng=` | - | 주변 장소 |
| GET | `/api/places/:id` | - | 장소 상세 |
| POST | `/api/places/:id/verify` | JWT | 위치 인증 |
| GET / POST | `/api/places/:id/guestbooks` | 위치 토큰 | 방명록 열람 / 작성 |
| GET | `/api/users/:id/stats` | JWT | 방문 장소 수, 글 수 |
| GET | `/api/users/:id/guestbooks/places` | JWT | 내가 다녀간 장소 |
| GET | `/api/users/:id/places/:placeId/guestbooks` | JWT | 특정 장소에 쓴 내 글 |
| POST | `/api/admin/places` | JWT + 관리자 | 장소 등록 |
| POST | `/api/admin/badges/:userId` | JWT + 관리자 | 뱃지 부여 |

## 로컬 실행

**Backend** (Go 1.26+, PostgreSQL + PostGIS)

```bash
cd backend
cp .env.example .env   # DATABASE_URL, JWT_SECRET(32자 이상), ADMIN_USER_IDS 입력
go run .
```

**Frontend** (Node 20.19+)

```bash
cd frontend
cp .env.example .env   # KAKAO_JS_KEY, KAKAO_REST_KEY 입력
npm install
npm run dev
```

## 보안 설계 메모

- 비밀번호는 bcrypt로 해시해서 저장하고, 응답 JSON에서는 `json:"-"`로 제외
- 처음에는 어드민 API를 공유 시크릿 헤더(`X-Admin-Secret`)로 보호했는데, 이 값이 프론트 번들에 포함돼 누구나 볼 수 있다는 문제가 있었습니다. **JWT 인증 + 서버 측 관리자 ID 화이트리스트(`ADMIN_USER_IDS`)** 방식으로 바꿔서, 클라이언트에는 비밀값이 남지 않게 했습니다.
- `JWT_SECRET`이 없거나 짧으면 기본값으로 뜨지 않고 기동 단계에서 실패하게 해서, 개발용 시크릿이 운영에 그대로 나가는 실수를 막았습니다.

<!-- TODO: 개발하면서 겪은 트러블슈팅 1~2개 (예: JWT 적용 시 토큰 누락 이슈, PostGIS 반경 검증) -->
