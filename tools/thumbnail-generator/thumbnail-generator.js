import sharp from 'sharp';
import fs from 'fs/promises';
import path from 'path';

export class ThumbnailGenerator {
  constructor() {
    this.defaultWidth = 1920;
    this.defaultHeight = 1080;
  }

  parseRatio(ratioString) {
    const parts = ratioString.split(':');
    if (parts.length !== 2) {
      throw new Error('비율은 "16:9" 형식으로 입력해주세요');
    }
    
    const width = parseInt(parts[0]);
    const height = parseInt(parts[1]);
    
    if (isNaN(width) || isNaN(height) || width <= 0 || height <= 0) {
      throw new Error('올바른 비율을 입력해주세요');
    }
    
    return { width, height };
  }

  parseColor(colorString) {
    if (colorString.startsWith('#')) {
      return colorString;
    }
    
    if (colorString.startsWith('rgb')) {
      return colorString;
    }
    
    const namedColors = {
      'red': '#ff0000',
      'green': '#00ff00',
      'blue': '#0000ff',
      'yellow': '#ffff00',
      'purple': '#800080',
      'orange': '#ffa500',
      'pink': '#ffc0cb',
      'brown': '#a52a2a',
      'gray': '#808080',
      'black': '#000000',
      'white': '#ffffff'
    };
    
    return namedColors[colorString.toLowerCase()] || '#ffffff';
  }

  calculateDimensions(ratio, maxWidth = 1920) {
    const aspectRatio = ratio.width / ratio.height;
    const width = maxWidth;
    const height = Math.round(width / aspectRatio);
    
    return { width, height };
  }

  async generateThumbnail(options) {
    const {
      logoPath,
      ratio = '16:9',
      backgroundColor = '#ffffff',
      outputPath = './output/thumbnail.png',
      size = 1920
    } = options;

    try {
      const parsedRatio = this.parseRatio(ratio);
      const bgColor = this.parseColor(backgroundColor);
      const dimensions = this.calculateDimensions(parsedRatio, size);

      await fs.access(logoPath);

      const logoBuffer = await fs.readFile(logoPath);
      const logoMetadata = await sharp(logoBuffer).metadata();

      const maxLogoSize = Math.min(dimensions.width, dimensions.height) * 0.4;
      
      let resizedLogo;
      if (logoMetadata.width > maxLogoSize || logoMetadata.height > maxLogoSize) {
        resizedLogo = await sharp(logoBuffer)
          .resize(Math.round(maxLogoSize), Math.round(maxLogoSize), {
            fit: 'inside',
            withoutEnlargement: false
          })
          .toBuffer();
      } else {
        resizedLogo = logoBuffer;
      }

      const resizedLogoMetadata = await sharp(resizedLogo).metadata();

      const thumbnail = await sharp({
        create: {
          width: dimensions.width,
          height: dimensions.height,
          channels: 4,
          background: bgColor
        }
      })
      .composite([{
        input: resizedLogo,
        left: Math.round((dimensions.width - resizedLogoMetadata.width) / 2),
        top: Math.round((dimensions.height - resizedLogoMetadata.height) / 2)
      }])
      .png()
      .toBuffer();

      const outputDir = path.dirname(outputPath);
      await fs.mkdir(outputDir, { recursive: true });
      
      await fs.writeFile(outputPath, thumbnail);

      return {
        success: true,
        outputPath,
        dimensions,
        logoSize: {
          width: resizedLogoMetadata.width,
          height: resizedLogoMetadata.height
        }
      };

    } catch (error) {
      throw new Error(`썸네일 생성 실패: ${error.message}`);
    }
  }
}