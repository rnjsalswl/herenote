(function(){
// HERE NOTE — mock data
// All Korean content for the location-based guestbook SNS.

const ME = {
  id: 'u_me',
  nickname: '재인',
  handle: '@here.jaein',
  bio: '발자국마다 한 줄씩. 서울 → 어디든.',
  stats: { visited: 18, notes: 42, stamps: 9 },
};

// Stamp / badge collection — earned by verifying at places
const STAMPS = [
  { id: 's1', name: '첫 발자국', icon: 'foot', earned: true, desc: '첫 장소 인증' },
  { id: 's2', name: '새벽 방문자', icon: 'moon', earned: true, desc: '오전 5시 이전 인증' },
  { id: 's3', name: '단골', icon: 'repeat', earned: true, desc: '같은 곳 5회 방문' },
  { id: 's4', name: '산 정상', icon: 'peak', earned: true, desc: '고도 500m 이상 인증' },
  { id: 's5', name: '바다 끝', icon: 'wave', earned: true, desc: '해안 장소 인증' },
  { id: 's6', name: '10곳 탐험', icon: 'compass', earned: true, desc: '서로 다른 10곳 방문' },
  { id: 's7', name: '글쟁이', icon: 'pen', earned: true, desc: '방명록 25개 작성' },
  { id: 's8', name: '공감 100', icon: 'heart', earned: true, desc: '받은 공감 100개' },
  { id: 's9', name: '터줏대감', icon: 'flag', earned: true, desc: '한 장소 첫 방문자' },
  { id: 's10', name: '심야 산책', icon: 'star', earned: false, desc: '자정 이후 인증' },
  { id: 's11', name: '전국 일주', icon: 'globe', earned: false, desc: '5개 도시 인증' },
  { id: 's12', name: '사진가', icon: 'camera', earned: false, desc: '사진 50장 첨부' },
];

// Places
const PLACES = [
  {
    id: 'p1',
    name: '경복궁 향원정',
    area: '서울 종로구',
    description: '연못 위 정자. 가을 단풍이 물에 비치는 자리.',
    radius: 100,
    distance: 40,
    unlocked: true,
    coord: { lat: 37.5811, lng: 126.9770 },
    visitors: 1284,
    notes: 312,
    cover: 'pond',
  },
  {
    id: 'p2',
    name: '연남동 경의선숲길',
    area: '서울 마포구',
    description: '기찻길 위에 깔린 잔디밭. 주말이면 돗자리가 빼곡.',
    radius: 200,
    distance: 120,
    unlocked: true,
    coord: { lat: 37.5602, lng: 126.9255 },
    visitors: 980,
    notes: 207,
    cover: 'park',
  },
  {
    id: 'p3',
    name: '북악산 팔각정',
    area: '서울 종로구',
    description: '서울이 한눈에. 야경 명소.',
    radius: 500,
    distance: 2400,
    unlocked: false,
    coord: { lat: 37.5996, lng: 126.9869 },
    visitors: 642,
    notes: 158,
    cover: 'peak',
  },
  {
    id: 'p4',
    name: '을지로 골목 노가리집',
    area: '서울 중구',
    description: '플라스틱 의자와 노란 전구. 퇴근길의 성지.',
    radius: 50,
    distance: 5100,
    unlocked: false,
    coord: { lat: 37.5662, lng: 126.9913 },
    visitors: 421,
    notes: 96,
    cover: 'alley',
  },
];

// Guestbook entries keyed by place id
const GUESTBOOKS = {
  p1: [
    {
      id: 'g1', userId: 'u2', nickname: '도윤', badge: 'GOLD', pinned: true,
      time: '2시간 전', likes: 48, liked: false, photo: 'pond',
      content: '비 갠 뒤 향원정. 물 냄새가 좋아서 한참을 앉아 있었다. 다음엔 단풍 들 때 다시 와야지.',
    },
    {
      id: 'g2', userId: 'u3', nickname: '서연', badge: 'BLUE', pinned: false,
      time: '5시간 전', likes: 21, liked: true, photo: null,
      content: '여기서 프로포즈 받았어요. 5년 전 오늘. 매년 와서 한 줄씩 남기는 중 🙂',
    },
    {
      id: 'g3', userId: 'u4', nickname: '민', badge: null, pinned: false,
      time: '어제', likes: 9, liked: false, photo: null,
      content: '평일 오전이라 사람이 거의 없었다. 새 소리랑 바람 소리뿐. 추천.',
    },
    {
      id: 'g4', userId: 'u5', nickname: '하루', badge: 'SPECIAL', pinned: false,
      time: '3일 전', likes: 64, liked: false, photo: 'park',
      content: '향원정 첫 방문자 뱃지 받음! 아무도 안 남긴 새 장소를 찾는 재미가 있다.',
    },
  ],
  p2: [
    {
      id: 'g5', userId: 'u6', nickname: '준호', badge: 'BLUE', pinned: true,
      time: '1시간 전', likes: 33, liked: false, photo: 'park',
      content: '돗자리 깔고 맥주 한 캔. 기차가 다니던 길이라니 신기하다.',
    },
    {
      id: 'g6', userId: 'u7', nickname: '유나', badge: null, pinned: false,
      time: '4시간 전', likes: 12, liked: false, photo: null,
      content: '강아지랑 산책하기 딱 좋아요. 그늘도 많고.',
    },
  ],
};

// Activity feed for Home (recent notes from places you've been)
const FEED = [
  {
    id: 'f1', place: '경복궁 향원정', area: '서울 종로구', placeId: 'p1',
    nickname: '도윤', badge: 'GOLD', time: '2시간 전', likes: 48, liked: false, photo: 'pond',
    content: '비 갠 뒤 향원정. 물 냄새가 좋아서 한참을 앉아 있었다.',
  },
  {
    id: 'f2', place: '연남동 경의선숲길', area: '서울 마포구', placeId: 'p2',
    nickname: '준호', badge: 'BLUE', time: '1시간 전', likes: 33, liked: false, photo: 'park',
    content: '돗자리 깔고 맥주 한 캔. 기차가 다니던 길이라니 신기하다.',
  },
  {
    id: 'f3', place: '경복궁 향원정', area: '서울 종로구', placeId: 'p1',
    nickname: '서연', badge: 'BLUE', time: '5시간 전', likes: 21, liked: true, photo: null,
    content: '여기서 프로포즈 받았어요. 매년 와서 한 줄씩 남기는 중 🙂',
  },
];

window.HN_DATA = { ME, STAMPS, PLACES, GUESTBOOKS, FEED };

})();
