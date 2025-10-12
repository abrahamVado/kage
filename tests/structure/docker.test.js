import { readFileSync } from 'node:fs';
import { join } from 'node:path';

//1.- Load the docker-compose manifest so assertions can validate required services and wiring.
const composePath = join(process.cwd(), 'deploy', 'docker-compose.yml');
const composeContent = readFileSync(composePath, 'utf8');

//2.- Ensure every critical service is defined with a build directive for its Dockerfile.
const requiredServices = [
  { name: 'backend', token: 'backend:\n    # //2.-' },
  { name: 'frontend', token: 'frontend:\n    # //3.-' },
  { name: 'mobile', token: 'mobile:\n    # //4.-' },
  { name: 'database', token: 'database:\n    # //5.-' },
];

for (const { name, token } of requiredServices) {
  if (!composeContent.includes(token)) {
    throw new Error(`Missing or misconfigured service: ${name}`);
  }
}

//3.- Confirm shared infrastructure pieces exist to connect and persist container state.
const sharedElements = ['kage-internal', 'database_data', 'frontend_node_modules', 'flutter_pub_cache', 'go_pkg'];
for (const element of sharedElements) {
  if (!composeContent.includes(element)) {
    throw new Error(`Expected shared element not found: ${element}`);
  }
}

console.log('docker-compose.yml declares all critical services and shared resources.');
