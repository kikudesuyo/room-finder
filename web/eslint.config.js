import eslint from '@eslint/js';
import prettier from 'eslint-config-prettier';
import tseslint from 'typescript-eslint';
import svelte from 'eslint-plugin-svelte';
import { defineConfig } from 'eslint/config';

export default defineConfig(
	{
		ignores: ['node_modules/**', '.svelte-kit/**', 'build/**']
	},
	eslint.configs.recommended,
	tseslint.configs.recommended,
	svelte.configs.recommended,
	prettier,
	svelte.configs.prettier,
	{
		rules: {
			'no-undef': 'off'
		}
	},
	{
		files: ['**/*.svelte', '**/*.svelte.ts', '**/*.svelte.js'],
		languageOptions: {
			parserOptions: {
				projectService: true,
				extraFileExtensions: ['.svelte'],
				parser: tseslint.parser
			}
		}
	}
);
