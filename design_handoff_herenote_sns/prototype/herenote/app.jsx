(function(){
// HERE NOTE — app shell, routing, tweaks
const { buildVars } = window.HN_UI;
const { TabBar } = window.HN_C;
const { LoginScreen, HomeScreen, ExploreScreen } = window.HN_S1;
const { PlaceScreen, ComposeScreen } = window.HN_S2;
const { ProfileScreen } = window.HN_S3;
const { useState, useEffect } = React;

const TWEAK_DEFAULTS = /*EDITMODE-BEGIN*/{
  "mood": "clean",
  "dark": false,
  "headline": "auto"
}/*EDITMODE-END*/;

const SANS = "'Pretendard', -apple-system, system-ui, sans-serif";
const SERIF = "'Gowun Batang', 'Pretendard', serif";

function App() {
  const [t, setTweak] = useTweaks(TWEAK_DEFAULTS);
  const data = window.HN_DATA;

  const [scale, setScale] = useState(1);
  useEffect(() => {
    const fit = () => {
      const s = Math.min(1, (window.innerHeight - 40) / 874, (window.innerWidth - 24) / 402);
      setScale(s > 0 ? s : 1);
    };
    fit();
    window.addEventListener('resize', fit);
    return () => window.removeEventListener('resize', fit);
  }, []);

  const [authed, setAuthed] = useState(false);
  const [tab, setTab] = useState('home');
  const [stack, setStack] = useState([]); // {name, params}
  const [verified, setVerified] = useState(() => new Set(data.PLACES.filter((p) => p.unlocked).map((p) => p.id)));
  const [extraNotes, setExtraNotes] = useState({}); // placeId -> [notes]

  const nav = {
    go: (name, params) => setStack((s) => [...s, { name, params: params || {} }]),
    back: () => setStack((s) => s.slice(0, -1)),
    reset: (route) => setStack(route ? [route] : []),
  };

  const placeById = (id) => data.PLACES.find((p) => p.id === id) || data.PLACES[0];
  const notesFor = (id) => [...(extraNotes[id] || []), ...(data.GUESTBOOKS[id] || [])];

  function openCompose() {
    const target = data.PLACES.find((p) => verified.has(p.id)) || data.PLACES[0];
    nav.go('compose', { placeId: target.id });
  }
  function handlePost(placeId, note) {
    setExtraNotes((m) => ({ ...m, [placeId]: [note, ...(m[placeId] || [])] }));
    setStack([{ name: 'place', params: { id: placeId } }]);
  }

  // ── theme vars ──
  const vars = buildVars(t.mood, t.dark);
  if (t.headline === 'sans') vars['--font-display'] = SANS;
  if (t.headline === 'serif') vars['--font-display'] = SERIF;

  const current = stack.length ? stack[stack.length - 1] : null;
  const onTabs = authed && !current;

  let screen;
  if (!authed) {
    screen = <LoginScreen onEnter={() => setAuthed(true)} />;
  } else if (current && current.name === 'place') {
    const place = placeById(current.params.id);
    screen = (
      <PlaceScreen
        place={place}
        notes={notesFor(place.id)}
        verified={verified.has(place.id)}
        nav={nav}
        onVerify={(id) => setVerified((s) => new Set(s).add(id))}
        onCompose={() => nav.go('compose', { placeId: place.id })}
      />
    );
  } else if (current && current.name === 'compose') {
    const place = placeById(current.params.placeId);
    screen = <ComposeScreen place={place} onClose={() => nav.back()} onPost={(note) => handlePost(place.id, note)} />;
  } else if (tab === 'home') {
    screen = <HomeScreen data={data} nav={nav} />;
  } else if (tab === 'map') {
    screen = <ExploreScreen data={data} nav={nav} />;
  } else {
    screen = <ProfileScreen data={data} nav={nav} />;
  }

  return (
    <div style={{
      minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center',
      background: t.dark ? '#1a1a1c' : '#e9e9ec', padding: '20px 0', overflow: 'hidden',
    }}>
      <div style={{ transform: `scale(${scale})`, transformOrigin: 'center center' }}>
      <IOSDevice dark={t.dark}>
        <div style={{
          ...vars, height: '100%', position: 'relative',
          display: 'flex', flexDirection: 'column',
          background: 'var(--bg)', color: 'var(--ink)', fontFamily: 'var(--font-sans)',
          WebkitFontSmoothing: 'antialiased',
        }}>
          <div className="hn-scroll" style={{ flex: 1, overflowY: 'auto', WebkitOverflowScrolling: 'touch' }}>
            {screen}
          </div>
          {onTabs && <TabBar active={tab} onChange={(id) => { setTab(id); }} onCompose={openCompose} />}
        </div>
      </IOSDevice>
      </div>

      <TweaksPanel title="Tweaks">
        <TweakSection label="장소 무드" />
        <TweakRadio label="무드" value={t.mood} options={[{ value: 'clean', label: '클린' }, { value: 'paper', label: '페이퍼' }, { value: 'editorial', label: '에디토리얼' }]} onChange={(v) => setTweak('mood', v)} />
        <TweakToggle label="다크 모드" value={t.dark} onChange={(v) => setTweak('dark', v)} />
        <TweakSection label="타이포그래피" />
        <TweakRadio label="헤드라인" value={t.headline} options={[{ value: 'auto', label: '자동' }, { value: 'sans', label: '산세리프' }, { value: 'serif', label: '세리프' }]} onChange={(v) => setTweak('headline', v)} />
      </TweaksPanel>
    </div>
  );
}

ReactDOM.createRoot(document.getElementById('root')).render(<App />);

})();
