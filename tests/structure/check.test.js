import { readFileSync, accessSync, constants } from 'node:fs';
import { join } from 'node:path';

const docPaths = [
  ['backend', 'docs', 'overview', 'README.md'],
  ['frontend', 'docs', 'overview', 'README.md'],
  ['mobile', 'docs', 'overview', 'README.md'],
  ['database', 'docs', 'overview', 'README.md'],
  ['deploy', 'docs', 'overview', 'README.md'],
  ['docs', 'handbook', 'README.md'],
  ['docs', 'system-overview', 'README.md']
];

//1.- Ensure all documentation files exist so teams can onboard consistently.
for (const segments of docPaths) {
  const filePath = join(process.cwd(), ...segments);
  try {
    //2.- Verify the file is accessible and contains non-empty content.
    accessSync(filePath, constants.R_OK);
    const content = readFileSync(filePath, 'utf8').trim();
    if (content.length === 0) {
      throw new Error(`File ${filePath} is empty`);
    }
  } catch (error) {
    console.error(`Missing or unreadable documentation: ${segments.join('/')}`);
    throw error;
  }
}

console.log('All critical documentation files are present and populated.');
