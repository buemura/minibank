interface GopherLogoProps {
  className?: string;
}

export default function GopherLogo({ className = "h-8 w-8" }: GopherLogoProps) {
  return (
    <img
      src="/icons/gopher.svg"
      alt="Gopher logo"
      className={className}
    />
  );
}
