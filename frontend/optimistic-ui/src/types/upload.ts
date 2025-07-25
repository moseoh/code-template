export interface ImageUpload {
  id: string;
  file: File;
  status: 'uploading' | 'completed' | 'error';
  uploadStartTime: Date;
  uploadCompleteTime?: Date;
  thumbnailStartTime?: Date;
  thumbnailCompleteTime?: Date;
  imageUrl?: string;
  thumbnailUrl?: string;
  error?: string;
}

export interface UploadResponse {
  success: boolean;
  imageUrl?: string;
  error?: string;
}