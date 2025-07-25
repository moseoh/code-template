import { ImageUpload } from '@/types/upload';
import ImageCard from './ImageCard';

interface ImageGridProps {
  uploads: ImageUpload[];
  onRemoveUpload?: (id: string) => void;
}

export default function ImageGrid({ uploads, onRemoveUpload }: ImageGridProps) {
  if (uploads.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">🖼️</div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">
          아직 업로드된 이미지가 없습니다
        </h3>
        <p className="text-gray-600">
          위에서 이미지를 업로드해보세요!
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      {uploads.map((upload) => (
        <ImageCard
          key={upload.id}
          upload={upload}
          onRemove={onRemoveUpload}
        />
      ))}
    </div>
  );
}