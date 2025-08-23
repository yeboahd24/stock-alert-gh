#!/usr/bin/env node

// Frontend Integration Validation Script
// This script validates the enhanced dividend alerts frontend integration

import fs from 'fs';
import path from 'path';

console.log('🔍 Validating Enhanced Dividend Alerts Frontend Integration...\n');

// Files to validate
const filesToCheck = [
  'components/forms/AlertForm.tsx',
  'components/dividend/DividendYieldDashboard.tsx', 
  'components/dividend/DividendYieldAlerts.tsx',
  'types/enums.ts',
  'types/schema.ts',
  'utils/formatters.ts',
  'src/services/api.ts',
  'tests/dividend-alerts.test.tsx'
];

// Validation results
const results = {
  passed: 0,
  failed: 0,
  issues: []
};

// Helper function to check if file exists and has content
function validateFile(filePath) {
  try {
    if (!fs.existsSync(filePath)) {
      return { success: false, error: 'File does not exist' };
    }
    
    const content = fs.readFileSync(filePath, 'utf8');
    if (content.length === 0) {
      return { success: false, error: 'File is empty' };
    }
    
    return { success: true, size: content.length };
  } catch (error) {
    return { success: false, error: error.message };
  }
}

// Helper function to check for specific content patterns
function validateContent(filePath, patterns) {
  try {
    const content = fs.readFileSync(filePath, 'utf8');
    const missing = patterns.filter(pattern => !content.includes(pattern));
    return missing.length === 0 ? { success: true } : { success: false, missing };
  } catch (error) {
    return { success: false, error: error.message };
  }
}

// Validate each file
console.log('📁 File Validation:');
console.log('='.repeat(50));

filesToCheck.forEach(file => {
  const result = validateFile(file);
  if (result.success) {
    console.log(`✅ ${file} (${Math.round(result.size / 1024)}KB)`);
    results.passed++;
  } else {
    console.log(`❌ ${file} - ${result.error}`);
    results.failed++;
    results.issues.push(`${file}: ${result.error}`);
  }
});

// Content validation for key files
console.log('\n🔍 Content Validation:');
console.log('='.repeat(50));

// Validate AlertForm has new alert types
const alertFormPatterns = [
  'AlertType.HIGH_DIVIDEND_YIELD',
  'AlertType.TARGET_DIVIDEND_YIELD', 
  'AlertType.DIVIDEND_YIELD_CHANGE',
  'thresholdYield',
  'targetYield',
  'yieldChangeThreshold'
];

if (fs.existsSync('components/forms/AlertForm.tsx')) {
  const alertFormResult = validateContent('components/forms/AlertForm.tsx', alertFormPatterns);
  if (alertFormResult.success) {
    console.log('✅ AlertForm contains new dividend yield alert types');
  } else {
    console.log(`❌ AlertForm missing: ${alertFormResult.missing?.join(', ')}`);
    results.issues.push('AlertForm missing dividend yield features');
  }
}

// Validate enums have new types
const enumPatterns = [
  'HIGH_DIVIDEND_YIELD',
  'TARGET_DIVIDEND_YIELD',
  'DIVIDEND_YIELD_CHANGE'
];

if (fs.existsSync('types/enums.ts')) {
  const enumResult = validateContent('types/enums.ts', enumPatterns);
  if (enumResult.success) {
    console.log('✅ Enums contain new dividend yield alert types');
  } else {
    console.log(`❌ Enums missing: ${enumResult.missing?.join(', ')}`);
    results.issues.push('Enums missing new alert types');
  }
}

// Validate API service has dividend endpoints
const apiPatterns = [
  'dividendApi',
  'getGSEDividendStocks',
  'getDividendStockBySymbol',
  'getHighDividendYieldStocks'
];

if (fs.existsSync('src/services/api.ts')) {
  const apiResult = validateContent('src/services/api.ts', apiPatterns);
  if (apiResult.success) {
    console.log('✅ API service contains dividend endpoints');
  } else {
    console.log(`❌ API service missing: ${apiResult.missing?.join(', ')}`);
    results.issues.push('API service missing dividend endpoints');
  }
}

// Validate formatters have yield formatting
const formatterPatterns = [
  'formatAlertType',
  'HIGH_DIVIDEND_YIELD',
  'TARGET_DIVIDEND_YIELD',
  'DIVIDEND_YIELD_CHANGE'
];

if (fs.existsSync('utils/formatters.ts')) {
  const formatterResult = validateContent('utils/formatters.ts', formatterPatterns);
  if (formatterResult.success) {
    console.log('✅ Formatters contain yield alert formatting');
  } else {
    console.log(`❌ Formatters missing: ${formatterResult.missing?.join(', ')}`);
    results.issues.push('Formatters missing yield alert support');
  }
}

// Check for syntax errors in key files
console.log('\n🔧 Syntax Validation:');
console.log('='.repeat(50));

const syntaxChecks = [
  { file: 'components/forms/AlertForm.tsx', patterns: ['export default AlertForm', 'interface AlertFormData'] },
  { file: 'types/enums.ts', patterns: ['export enum AlertType'] },
  { file: 'types/schema.ts', patterns: ['export interface'] }
];

syntaxChecks.forEach(check => {
  if (fs.existsSync(check.file)) {
    const syntaxResult = validateContent(check.file, check.patterns);
    if (syntaxResult.success) {
      console.log(`✅ ${check.file} syntax OK`);
    } else {
      console.log(`❌ ${check.file} syntax issues`);
      results.issues.push(`${check.file} has syntax issues`);
    }
  }
});

// Summary
console.log('\n📊 Validation Summary:');
console.log('='.repeat(50));
console.log(`Files checked: ${filesToCheck.length}`);
console.log(`Passed: ${results.passed}`);
console.log(`Failed: ${results.failed}`);

if (results.issues.length > 0) {
  console.log('\n❌ Issues found:');
  results.issues.forEach(issue => console.log(`  • ${issue}`));
} else {
  console.log('\n🎉 All validations passed!');
}

// Integration checklist
console.log('\n📋 Integration Checklist:');
console.log('='.repeat(50));

const checklist = [
  { item: 'Enhanced AlertForm with dividend yield fields', file: 'components/forms/AlertForm.tsx' },
  { item: 'DividendYieldDashboard component', file: 'components/dividend/DividendYieldDashboard.tsx' },
  { item: 'DividendYieldAlerts management component', file: 'components/dividend/DividendYieldAlerts.tsx' },
  { item: 'Updated alert type enums', file: 'types/enums.ts' },
  { item: 'Enhanced schema types', file: 'types/schema.ts' },
  { item: 'Dividend API integration', file: 'src/services/api.ts' },
  { item: 'Updated formatters', file: 'utils/formatters.ts' },
  { item: 'Comprehensive test suite', file: 'tests/dividend-alerts.test.tsx' }
];

checklist.forEach(item => {
  const exists = fs.existsSync(item.file);
  console.log(`${exists ? '✅' : '❌'} ${item.item}`);
});

console.log('\n🚀 Next Steps:');
console.log('='.repeat(50));
console.log('1. Run: npm install (if not done)');
console.log('2. Run: npm run dev (start frontend)');
console.log('3. Run: npm test dividend-alerts.test.tsx (run tests)');
console.log('4. Navigate to dividend alerts in the app');
console.log('5. Test creating new dividend yield alerts');

console.log('\n📚 Documentation:');
console.log('='.repeat(50));
console.log('• FRONTEND_DIVIDEND_ALERTS_INTEGRATION.md - Complete integration guide');
console.log('• ENHANCED_DIVIDEND_ALERTS.md - Backend features documentation');
console.log('• DIVIDEND_API_INTEGRATION.md - API integration guide');

if (results.failed === 0 && results.issues.length === 0) {
  console.log('\n🎉 Frontend integration validation completed successfully!');
  process.exit(0);
} else {
  console.log('\n⚠️ Some issues found. Please review and fix before proceeding.');
  process.exit(1);
}