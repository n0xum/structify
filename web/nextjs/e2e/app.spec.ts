import { test, expect } from '@playwright/test';

test.describe('Structify App E2E', () => {
  test('generates repository implementation from interface', async ({ page }) => {
    await page.goto('/');

    // Ensure page is loaded
    await expect(page.locator('text=structify').first()).toBeVisible();

    // Select Interface Repository mode
    await page.getByRole('tab', { name: 'Interface Repository' }).click();

    // Set the package name
    const pkgInput = page.getByLabel('Package name');
    await pkgInput.fill('myrepo');

    // The editor requires complex interaction, let's paste into it or mock the textarea if it uses native textarea.
    // The editor in UIW React CodeMirror uses a contenteditable element.
    const editor = page.locator('.cm-content');
    
    // Clear editor and enter new code. Using clear and type.
    const code = "package models\nimport \"context\"\n\ntype Product struct {\n  ID int64 `db:\"pk\"`\n  Name string\n}\n\ntype ProductRepository interface {\n  FindByID(ctx context.Context, id int64) (*Product, error)\n}";
    
    // CodeMirror handling
    await editor.click();
    
    // Clear the existing text (Ctrl+A / Cmd+A, then Delete)
    const modifier = process.platform === 'darwin' ? 'Meta' : 'Control';
    await page.keyboard.press(`${modifier}+A`);
    await page.keyboard.press('Backspace');

    // Insert text
    await page.keyboard.insertText(code);

    // Click Generate
    await page.getByRole('button', { name: 'Generate output (Ctrl+Enter)' }).click();

    // Verify output
    // Wait for the generated text to appear in the output panel.
    const outputPanel = page.locator('pre code');
    await expect(outputPanel).toContainText('type ProductRepositoryImpl struct', { timeout: 10000 });
    await expect(outputPanel).toContainText('package myrepo');
    await expect(outputPanel).toContainText('func (r *ProductRepositoryImpl) FindByID');
  });

  test('generates SQL schema', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('tab', { name: 'SQL Schema' }).click();

    const editor = page.locator('.cm-content');
    const code = "package models\ntype User struct {\n  ID int64 `db:\"pk\"`\n}";

    await editor.click();
    const modifier = process.platform === 'darwin' ? 'Meta' : 'Control';
    await page.keyboard.press(`${modifier}+A`);
    await page.keyboard.press('Backspace');
    await page.keyboard.insertText(code);

    await page.getByRole('button', { name: 'Generate output (Ctrl+Enter)' }).click();

    const outputPanel = page.locator('pre code');
    await expect(outputPanel).toContainText('CREATE TABLE "user"', { timeout: 10000 });
  });
});
