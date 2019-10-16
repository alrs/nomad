/* eslint-env node */
// FIXME Vault has an entry in .eslintignore to skip Storybook altogether…???
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text } from '@storybook/addon-knobs';
import notes from './copy-button.md';


storiesOf('CopyButton/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add('CopyButton', () => ({
    template: hbs`
      <h5 class="title is-5">Copy Button</h5>
      <span class="tag is-hollow is-small no-text-transform">
        {{clipboardText}}
        {{copy-button clipboardText=clipboardText}}
      </span>
    `,
    context: {
      clipboardText: text('Clipboard Text', 'e8c898a0-794b-9063-7a7f-bf0c4a405f83'),
    },
  }),
  {notes}
  );
