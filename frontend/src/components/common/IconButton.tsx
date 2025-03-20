'use client';

type IconButtonProps = {
  icon: React.ReactNode;
  onClick: () => void;
  className?: string;
  label?: string;
  children?: React.ReactNode;
  type?: 'button' | 'submit';
};

export function IconButton({ 
  icon, 
  onClick, 
  className = '', 
  label,
  children,
  type = 'button'
}: IconButtonProps) {
  return (
    <button
      type={type}
      className={`inline-flex items-center ${className}`}
      onClick={onClick}
      aria-label={label}
    >
      {icon}
      {children}
    </button>
  );
} 