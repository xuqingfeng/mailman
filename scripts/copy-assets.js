#!/usr/bin/env node
'use strict';

// Vendors the front-end libraries from node_modules into
// ui/assets/lib/<name>/<version>, which is what the Go binary serves.
// The version is read from each installed package, so package.json stays the
// single source of truth for what ships. Older versions of a managed package
// are pruned.

const fs = require('fs');
const path = require('path');

const projectRoot = path.resolve(__dirname, '..');
const nodeModulesDir = path.join(projectRoot, 'node_modules');
const libDir = path.join(projectRoot, 'ui', 'assets', 'lib');

// `from` is relative to the installed package; `to` is relative to
// ui/assets/lib/<name>/<version>.
const assets = [
  {
    package: 'bootstrap',
    name: 'bootstrap',
    entries: [
      ['dist/css/bootstrap.min.css', 'css/bootstrap.min.css'],
      ['dist/css/bootstrap.min.css.map', 'css/bootstrap.min.css.map'],
      ['dist/js/bootstrap.min.js', 'js/bootstrap.min.js'],
      ['dist/js/bootstrap.min.js.map', 'js/bootstrap.min.js.map'],
    ],
  },
  {
    package: '@selectize/selectize',
    name: 'selectize',
    // copy the whole dist/ tree, stripping the dist/ prefix
    entries: [['dist', '.']],
  },
  {
    package: 'jquery',
    name: 'jquery',
    entries: [
      ['dist/jquery.min.js', 'jquery.min.js'],
      ['dist/jquery.min.map', 'jquery.min.map'],
    ],
  },
];

function installedVersion(pkg) {
  const manifest = path.join(nodeModulesDir, pkg, 'package.json');
  if (!fs.existsSync(manifest)) {
    throw new Error(`${pkg} is not installed, run "npm install" first`);
  }
  return JSON.parse(fs.readFileSync(manifest, 'utf8')).version;
}

function copyEntry(source, target) {
  if (!fs.existsSync(source)) {
    throw new Error(`missing asset: ${source}`);
  }
  if (fs.statSync(source).isDirectory()) {
    fs.mkdirSync(target, { recursive: true });
    for (const name of fs.readdirSync(source)) {
      fs.cpSync(path.join(source, name), path.join(target, name), { recursive: true });
    }
  } else {
    fs.mkdirSync(path.dirname(target), { recursive: true });
    fs.copyFileSync(source, target);
  }
}

function pruneVersions(asset, currentVersion) {
  const packageDir = path.join(libDir, asset.name);
  if (!fs.existsSync(packageDir)) {
    return;
  }
  for (const name of fs.readdirSync(packageDir)) {
    if (name === currentVersion) {
      continue;
    }
    const stale = path.join(packageDir, name);
    if (!fs.statSync(stale).isDirectory()) {
      continue;
    }
    fs.rmSync(stale, { recursive: true, force: true });
    console.log(`pruned  ${path.relative(projectRoot, stale)}`);
  }
}

function main() {
  for (const asset of assets) {
    const version = installedVersion(asset.package);
    const packageRoot = path.join(nodeModulesDir, asset.package);
    const targetRoot = path.join(libDir, asset.name, version);

    fs.rmSync(targetRoot, { recursive: true, force: true });
    fs.mkdirSync(targetRoot, { recursive: true });

    for (const [from, to] of asset.entries) {
      copyEntry(path.join(packageRoot, from), path.join(targetRoot, to));
    }

    console.log(`copied  ${asset.package}@${version} -> ${path.relative(projectRoot, targetRoot)}`);
    pruneVersions(asset, version);
  }
}

main();
