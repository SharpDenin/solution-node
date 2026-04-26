import { useState, useRef } from 'react';

interface Props {
  onUpload: (file: File) => void;
  onRemove: () => void;
  imageUrl?: string;
}

export const ImageUploader: React.FC<Props> = ({ onUpload, onRemove, imageUrl }) => {
  const [dragActive, setDragActive] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    const files = e.dataTransfer.files;
    if (files && files[0]) {
      onUpload(files[0]);
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) onUpload(file);
  };

  return (
    <div>
      <div
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
        style={{
          border: `2px dashed ${dragActive ? '#16a34a' : '#ccc'}`,
          borderRadius: 8,
          padding: '16px',
          textAlign: 'center' as const,
          cursor: 'pointer',
          transition: 'all 0.2s',
          backgroundColor: dragActive ? '#f0fdf4' : '#fafafa',
          marginBottom: 8,
        }}
        onClick={() => fileInputRef.current?.click()}
      >
        <input
          type="file"
          accept="image/*"
          onChange={handleFileSelect}
          style={{ display: 'none' }}
          ref={fileInputRef}
        />
        <div style={{ color: '#555' }}>📸 Нажмите или перетащите изображение</div>
      </div>
      {imageUrl && (
        <div style={{ position: 'relative' as const, display: 'inline-block' }}>
          <img src={imageUrl} style={{ maxWidth: 200, maxHeight: 150, borderRadius: 8 }} alt="preview" />
          <button
            onClick={(e) => { e.stopPropagation(); onRemove(); }}
            style={{
              position: 'absolute' as const,
              top: -8,
              right: -8,
              background: '#ef4444',
              color: 'white',
              border: 'none',
              borderRadius: '50%',
              width: 24,
              height: 24,
              cursor: 'pointer',
              fontSize: 12,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            ✕
          </button>
        </div>
      )}
    </div>
  );
};