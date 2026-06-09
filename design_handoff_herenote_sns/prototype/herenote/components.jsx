(function(){
// HERE NOTE — shared UI components
const { Icon, Badge, Avatar, MapThumb, PhotoBlock } = window.HN_UI;
const { useState } = React;

// ── Pill ──────────────────────────────────────────────────────
function Pill({ children, solid = false, mono = false, style = {} }) {
  return (
    <span style={{
      display: 'inline-flex', alignItems: 'center', gap: 5,
      padding: '4px 10px', borderRadius: 999, fontSize: 12, lineHeight: 1,
      fontFamily: mono ? 'var(--font-mono)' : 'var(--font-sans)',
      fontWeight: 600, letterSpacing: mono ? '0.02em' : '0',
      background: solid ? 'var(--ink)' : 'var(--surface-2)',
      color: solid ? 'var(--inv-ink)' : 'var(--ink-2)',
      border: solid ? 'none' : '1px solid var(--line)',
      ...style,
    }}>{children}</span>
  );
}

// ── Top bar ───────────────────────────────────────────────────
function TopBar({ title, big, onBack, right, transparent = false, mono = false }) {
  return (
    <div style={{
      position: 'sticky', top: 0, zIndex: 30,
      paddingTop: 54, // clears status bar / dynamic island
      background: transparent ? 'transparent' : 'color-mix(in srgb, var(--bg) 82%, transparent)',
      backdropFilter: transparent ? 'none' : 'blur(18px) saturate(150%)',
      WebkitBackdropFilter: transparent ? 'none' : 'blur(18px) saturate(150%)',
      borderBottom: transparent ? '1px solid transparent' : '1px solid var(--line)',
    }}>
      <div style={{
        height: 50, display: 'flex', alignItems: 'center', gap: 6,
        padding: '0 10px',
      }}>
        {onBack && (
          <button onClick={onBack} style={iconBtn}>
            <Icon name="chevL" size={24} stroke={2.4} />
          </button>
        )}
        <div style={{
          flex: 1, minWidth: 0, display: 'flex', alignItems: 'center',
          justifyContent: onBack ? 'center' : 'flex-start', paddingLeft: onBack ? 0 : 8,
        }}>
          {big ? (
            <span style={{
              fontFamily: 'var(--font-display)', fontWeight: 'var(--disp-weight)',
              fontSize: 21, letterSpacing: mono ? '0.14em' : 'var(--disp-ls)',
              fontVariant: mono ? 'all-small-caps' : 'normal',
              color: 'var(--ink)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis',
            }}>{title}</span>
          ) : (
            <span style={{
              fontFamily: 'var(--font-sans)', fontWeight: 700, fontSize: 16,
              color: 'var(--ink)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis',
            }}>{title}</span>
          )}
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 2, minWidth: onBack ? 40 : 'auto' }}>
          {right}
        </div>
      </div>
    </div>
  );
}

const iconBtn = {
  width: 40, height: 40, borderRadius: '50%', border: 'none', background: 'transparent',
  color: 'var(--ink)', display: 'flex', alignItems: 'center', justifyContent: 'center',
  cursor: 'pointer', flexShrink: 0, WebkitTapHighlightColor: 'transparent',
};

// ── Bottom tab bar ────────────────────────────────────────────
function TabBar({ active, onChange, onCompose }) {
  const tab = (id, icon, label) => {
    const on = active === id;
    return (
      <button onClick={() => onChange(id)} style={{
        flex: 1, border: 'none', background: 'transparent', cursor: 'pointer',
        display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 3,
        color: on ? 'var(--ink)' : 'var(--ink-3)', paddingTop: 2,
        WebkitTapHighlightColor: 'transparent',
      }}>
        <Icon name={icon} size={25} stroke={on ? 2.4 : 1.9} fill={on && icon === 'home'} />
        <span style={{ fontSize: 10.5, fontWeight: on ? 700 : 500, fontFamily: 'var(--font-sans)' }}>{label}</span>
      </button>
    );
  };
  return (
    <div style={{
      flexShrink: 0, zIndex: 40,
      background: 'color-mix(in srgb, var(--bg) 88%, transparent)',
      backdropFilter: 'blur(18px) saturate(150%)', WebkitBackdropFilter: 'blur(18px) saturate(150%)',
      borderTop: '1px solid var(--line)',
      display: 'flex', alignItems: 'center', gap: 4, padding: '8px 14px 26px',
    }}>
      {tab('home', 'home', '홈')}
      {tab('map', 'map', '지도')}
      <button onClick={onCompose} aria-label="방명록 쓰기" style={{
        width: 54, height: 42, borderRadius: 15, border: 'none', cursor: 'pointer',
        background: 'var(--ink)', color: 'var(--inv-ink)', flexShrink: 0,
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        boxShadow: 'var(--shadow)', WebkitTapHighlightColor: 'transparent',
      }}>
        <Icon name="plus" size={24} stroke={2.6} />
      </button>
      {tab('profile', 'user', '프로필')}
    </div>
  );
}

// ── Note card (feed + place feed) ─────────────────────────────
function NoteCard({ note, place, onOpen, pinnable = true }) {
  const [liked, setLiked] = useState(!!note.liked);
  const [count, setCount] = useState(note.likes);
  const toggleLike = () => { setLiked((v) => !v); setCount((c) => c + (liked ? -1 : 1)); };
  return (
    <article style={{
      background: 'var(--card)', borderRadius: 'var(--radius-lg)',
      border: '1px solid var(--line)', overflow: 'hidden',
      boxShadow: 'var(--shadow)',
    }}>
      {/* header */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '12px 14px 10px' }}>
        <Avatar name={note.nickname} size={40} />
        <div style={{ flex: 1, minWidth: 0 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 5 }}>
            <span style={{ fontWeight: 700, fontSize: 14.5, color: 'var(--ink)' }}>{note.nickname}</span>
            <Badge type={note.badge} />
            {note.pinned && pinnable && <Pill style={{ padding: '2px 7px', fontSize: 10.5 }}>고정</Pill>}
          </div>
          {place && (
            <button onClick={onOpen} style={{
              border: 'none', background: 'transparent', padding: 0, cursor: 'pointer',
              display: 'flex', alignItems: 'center', gap: 3, color: 'var(--ink-2)',
              fontSize: 12.5, fontWeight: 500, marginTop: 1,
            }}>
              <Icon name="pin" size={13} stroke={2.2} />
              <span style={{ fontWeight: 600, color: 'var(--ink-2)' }}>{place}</span>
            </button>
          )}
        </div>
        <span style={{ fontFamily: 'var(--font-mono)', fontSize: 11, color: 'var(--ink-3)', flexShrink: 0 }}>{note.time}</span>
      </div>

      {/* photo */}
      {note.photo && (
        <div style={{ padding: '0 14px 12px' }}>
          <PhotoBlock kind={note.photo} height={196} />
        </div>
      )}

      {/* content */}
      <p style={{
        padding: note.photo ? '0 16px 4px' : '0 16px 4px', margin: 0,
        fontSize: 15, lineHeight: 1.62, color: 'var(--ink)', whiteSpace: 'pre-wrap',
        wordBreak: 'keep-all',
      }}>{note.content}</p>

      {/* actions */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 4, padding: '10px 12px 12px' }}>
        <button onClick={toggleLike} style={{
          display: 'inline-flex', alignItems: 'center', gap: 6, border: 'none',
          background: 'transparent', cursor: 'pointer', padding: '6px 8px', borderRadius: 999,
          color: liked ? 'var(--ink)' : 'var(--ink-2)', WebkitTapHighlightColor: 'transparent',
        }}>
          <span style={{ display: 'inline-flex', transform: liked ? 'scale(1.18)' : 'scale(1)', transition: 'transform .18s cubic-bezier(.2,.9,.3,1.4)' }}>
            <Icon name="heart" size={20} stroke={2} fill={liked} />
          </span>
          <span style={{ fontSize: 13, fontWeight: 700, fontVariantNumeric: 'tabular-nums' }}>{count}</span>
        </button>
        <button style={{
          display: 'inline-flex', alignItems: 'center', gap: 6, border: 'none',
          background: 'transparent', cursor: 'pointer', padding: '6px 8px', borderRadius: 999,
          color: 'var(--ink-2)', WebkitTapHighlightColor: 'transparent',
        }}>
          <Icon name="comment" size={19} stroke={2} />
          <span style={{ fontSize: 13, fontWeight: 600 }}>공감</span>
        </button>
        <div style={{ flex: 1 }} />
        <button style={{ ...iconBtn, width: 34, height: 34, color: 'var(--ink-2)' }}>
          <Icon name="bookmark" size={18} stroke={2} />
        </button>
      </div>
    </article>
  );
}

window.HN_C = { Pill, TopBar, TabBar, NoteCard, iconBtn };

})();
