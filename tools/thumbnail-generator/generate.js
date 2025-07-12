#!/usr/bin/env node

import { AutoThumbnailGenerator } from './auto-generator.js';

const generator = new AutoThumbnailGenerator();
await generator.generateAll();