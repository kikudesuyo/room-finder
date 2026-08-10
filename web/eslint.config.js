import eslint from '@eslint/js';
import prettier from 'eslint-config-prettier';
import tseslint from 'typescript-eslint';
import svelte from 'eslint-plugin-svelte';

export default tseslint.config(
	eslint.configs.recommended,
	...tseslint.configs.recommended,
	...svelte.configs.recommended,
	prettier,
	{
		files: ['**/*.ts', '**/*.svelte'],
		languageOptions: {
			parserOptions: {
				projectService: true
			}
		}
	},
	{
		ignores: ['.svelte-kit/**', 'build/**', 'node_modules/**']
	}
);
