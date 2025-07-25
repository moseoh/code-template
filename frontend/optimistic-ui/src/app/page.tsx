'use client';

import { useImageUploads } from '@/hooks/useImageUploads';
import UploadArea from '@/components/UploadArea';
import ImageGrid from '@/components/ImageGrid';

export default function Home() {
  const { uploads, uploadMultipleFiles, removeUpload } = useImageUploads();

  const handleFilesSelected = (files: FileList) => {
    uploadMultipleFiles(files);
  };

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          {/* 헤더 */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              Optimistic UI 이미지 업로드
            </h1>
            <p className="text-gray-600">
              실시간으로 업로드 상태를 확인하며 이미지를 업로드해보세요
            </p>
          </div>

          {/* 업로드 영역 */}
          <div className="mb-8">
            <UploadArea onFilesSelected={handleFilesSelected} />
          </div>

          {/* 통계 정보 */}
          {uploads.length > 0 && (
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-8">
              <div className="flex flex-wrap gap-6 text-sm">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
                  <span>업로드 중: {uploads.filter(u => u.status === 'uploading').length}</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                  <span>완료: {uploads.filter(u => u.status === 'completed').length}</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                  <span>오류: {uploads.filter(u => u.status === 'error').length}</span>
                </div>
                <div className="ml-auto font-medium">
                  총 {uploads.length}개 파일
                </div>
              </div>
            </div>
          )}

          {/* 이미지 그리드 */}
          <ImageGrid uploads={uploads} onRemoveUpload={removeUpload} />
        </div>
      </div>
    </main>
  );
}
