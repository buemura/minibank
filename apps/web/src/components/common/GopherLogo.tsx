interface GopherLogoProps {
  className?: string;
}

export default function GopherLogo({ className = "h-8 w-8" }: GopherLogoProps) {
  return (
    <svg
      className={className}
      viewBox="0 0 200 200"
      xmlns="http://www.w3.org/2000/svg"
    >
      {/* Body */}
      <ellipse cx="100" cy="120" rx="55" ry="60" fill="#65C4E0" />
      {/* Head */}
      <circle cx="100" cy="70" r="45" fill="#65C4E0" />
      {/* Eyes (whites) */}
      <circle cx="82" cy="60" r="14" fill="white" />
      <circle cx="118" cy="60" r="14" fill="white" />
      {/* Eyes (pupils) */}
      <circle cx="85" cy="60" r="6" fill="#1A1A1A" />
      <circle cx="121" cy="60" r="6" fill="#1A1A1A" />
      {/* Eye shine */}
      <circle cx="87" cy="58" r="2" fill="white" />
      <circle cx="123" cy="58" r="2" fill="white" />
      {/* Nose */}
      <ellipse cx="100" cy="72" rx="4" ry="3" fill="#D4A76A" />
      {/* Mouth */}
      <path
        d="M 90 80 Q 100 90 110 80"
        stroke="#1A1A1A"
        strokeWidth="2"
        fill="none"
        strokeLinecap="round"
      />
      {/* Teeth */}
      <rect x="96" y="78" width="8" height="6" rx="1" fill="white" stroke="#1A1A1A" strokeWidth="0.5" />
      {/* Left ear */}
      <ellipse cx="62" cy="38" rx="8" ry="16" fill="#65C4E0" transform="rotate(-20, 62, 38)" />
      <ellipse cx="62" cy="38" rx="4" ry="10" fill="#4AB0D4" transform="rotate(-20, 62, 38)" />
      {/* Right ear */}
      <ellipse cx="138" cy="38" rx="8" ry="16" fill="#65C4E0" transform="rotate(20, 138, 38)" />
      <ellipse cx="138" cy="38" rx="4" ry="10" fill="#4AB0D4" transform="rotate(20, 138, 38)" />
      {/* Left arm */}
      <ellipse cx="50" cy="125" rx="12" ry="22" fill="#65C4E0" transform="rotate(15, 50, 125)" />
      {/* Right arm */}
      <ellipse cx="150" cy="125" rx="12" ry="22" fill="#65C4E0" transform="rotate(-15, 150, 125)" />
      {/* Belly */}
      <ellipse cx="100" cy="130" rx="30" ry="32" fill="#B8DFF0" />
      {/* Left foot */}
      <ellipse cx="80" cy="178" rx="16" ry="8" fill="#D4A76A" />
      {/* Right foot */}
      <ellipse cx="120" cy="178" rx="16" ry="8" fill="#D4A76A" />
    </svg>
  );
}
