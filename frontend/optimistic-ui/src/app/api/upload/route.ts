import { NextRequest, NextResponse } from 'next/server';
import { writeFile, mkdir } from 'fs/promises';
import path from 'path';

// 썸네일 생성 비동기 함수 (1초 고정)
async function generateThumbnailAsync(fileName: string) {
  try {
    await new Promise(resolve => setTimeout(resolve, 1000));
    console.log(`썸네일 생성 완료: ${fileName}`);
    // 여기서 실제 썸네일 생성 로직이 들어갈 예정
  } catch (error) {
    console.error(`썸네일 생성 실패: ${fileName}`, error);
  }
}

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData();
    const file = formData.get('file') as File;
    
    if (!file) {
      return NextResponse.json({ 
        success: false, 
        error: '파일이 없습니다.' 
      }, { status: 400 });
    }

    const bytes = await file.arrayBuffer();
    const buffer = Buffer.from(bytes);

    // 업로드 시뮬레이션을 위한 1초 고정 딜레이
    await new Promise(resolve => setTimeout(resolve, 1000));

    // 파일명 생성 (타임스탬프 + 원본 파일명)
    const timestamp = Date.now();
    const fileName = `${timestamp}_${file.name}`;
    const uploadDir = path.join(process.cwd(), 'public/uploads');
    
    // 업로드 디렉토리 생성
    try {
      await mkdir(uploadDir, { recursive: true });
    } catch (error) {
      // 디렉토리가 이미 존재하는 경우 무시
    }

    // 파일 저장
    const filePath = path.join(uploadDir, fileName);
    await writeFile(filePath, buffer);

    const imageUrl = `/uploads/${fileName}`;

    // 썸네일 생성을 비동기로 시작 (응답을 블록하지 않음)
    generateThumbnailAsync(fileName);

    // 업로드 완료 즉시 응답
    return NextResponse.json({
      success: true,
      imageUrl,
    });

  } catch (error) {
    console.error('파일 업로드 에러:', error);
    return NextResponse.json({ 
      success: false, 
      error: '파일 업로드에 실패했습니다.' 
    }, { status: 500 });
  }
}