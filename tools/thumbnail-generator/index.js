#!/usr/bin/env node

import { Command } from 'commander';
import chalk from 'chalk';
import { ThumbnailGenerator } from './thumbnail-generator.js';
import path from 'path';

const program = new Command();
const generator = new ThumbnailGenerator();

program
  .name('thumbnail-generator')
  .description('로고와 배경색을 사용한 썸네일 이미지 생성기')
  .version('1.0.0');

program
  .command('generate')
  .alias('gen')
  .description('썸네일 이미지 생성')
  .requiredOption('-l, --logo <path>', '로고 이미지 파일 경로 (SVG 또는 PNG)')
  .option('-r, --ratio <ratio>', '이미지 비율 (예: 16:9, 4:3, 1:1)', '16:9')
  .option('-bg, --background <color>', '배경색 (hex, rgb, 또는 색상명)', '#ffffff')
  .option('-o, --output <path>', '출력 파일 경로', './output/thumbnail.png')
  .option('-s, --size <width>', '이미지 너비 (높이는 비율에 따라 자동 계산)', '1920')
  .action(async (options) => {
    try {
      console.log(chalk.blue('📸 썸네일 생성을 시작합니다...'));
      console.log(chalk.gray(`로고: ${options.logo}`));
      console.log(chalk.gray(`비율: ${options.ratio}`));
      console.log(chalk.gray(`배경색: ${options.background}`));
      console.log(chalk.gray(`출력: ${options.output}`));
      
      const result = await generator.generateThumbnail({
        logoPath: options.logo,
        ratio: options.ratio,
        backgroundColor: options.background,
        outputPath: options.output,
        size: parseInt(options.size)
      });

      console.log(chalk.green('✅ 썸네일이 성공적으로 생성되었습니다!'));
      console.log(chalk.yellow(`📁 출력 파일: ${result.outputPath}`));
      console.log(chalk.yellow(`📐 크기: ${result.dimensions.width}x${result.dimensions.height}`));
      console.log(chalk.yellow(`🖼️  로고 크기: ${result.logoSize.width}x${result.logoSize.height}`));
      
    } catch (error) {
      console.error(chalk.red('❌ 오류:'), error.message);
      process.exit(1);
    }
  });

program
  .command('colors')
  .description('사용 가능한 색상 목록 표시')
  .action(() => {
    console.log(chalk.blue('🎨 사용 가능한 색상:'));
    console.log('');
    console.log(chalk.red('● red') + '     ' + chalk.green('● green') + '   ' + chalk.blue('● blue'));
    console.log(chalk.yellow('● yellow') + '  ' + chalk.magenta('● purple') + '  ' + '● orange');
    console.log('● pink     ● brown    ● gray');
    console.log('● black    ● white');
    console.log('');
    console.log('또는 hex 코드 (#ff0000) 또는 rgb 값 (rgb(255,0,0))을 사용할 수 있습니다.');
  });

program
  .command('examples')
  .description('사용 예제 표시')
  .action(() => {
    console.log(chalk.blue('📖 사용 예제:'));
    console.log('');
    console.log(chalk.gray('기본 16:9 썸네일 생성:'));
    console.log('  ' + chalk.cyan('thumbnail-generator generate -l ./assets/logo.png'));
    console.log('');
    console.log(chalk.gray('4:3 비율, 빨간 배경:'));
    console.log('  ' + chalk.cyan('thumbnail-generator generate -l ./assets/logo.svg -r 4:3 -bg red'));
    console.log('');
    console.log(chalk.gray('사용자 정의 hex 색상:'));
    console.log('  ' + chalk.cyan('thumbnail-generator generate -l ./assets/logo.png -bg "#ff6b6b"'));
    console.log('');
    console.log(chalk.gray('사용자 정의 출력 경로:'));
    console.log('  ' + chalk.cyan('thumbnail-generator generate -l ./assets/logo.png -o ./my-thumbnail.png'));
    console.log('');
    console.log(chalk.gray('큰 크기 (4K):'));
    console.log('  ' + chalk.cyan('thumbnail-generator generate -l ./assets/logo.png -s 3840'));
  });

if (process.argv.length === 2) {
  program.help();
}

program.parse();