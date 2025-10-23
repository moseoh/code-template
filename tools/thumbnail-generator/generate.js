#!/usr/bin/env node

import { Command } from 'commander';
import { AutoThumbnailGenerator } from './auto-generator.js';

const program = new Command();

program
  .name('thumbnail-generator')
  .description('로고와 배경색을 사용한 썸네일 이미지 생성기')
  .option('-c, --config <path>', 'config 파일 경로 지정', './config.json')
  .parse();

const options = program.opts();

const generator = new AutoThumbnailGenerator(options.config);
await generator.generateAll();