// Populate the sidebar
//
// This is a script, and not included directly in the page, to control the total size of the book.
// The TOC contains an entry for each page, so if each page includes a copy of the TOC,
// the total size of the page becomes O(n**2).
class MDBookSidebarScrollbox extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this.innerHTML = '<ol class="chapter"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace.html"><strong aria-hidden="true">1.</strong> carapace</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen.html"><strong aria-hidden="true">1.1.</strong> Gen</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/dashAnyCompletion.html"><strong aria-hidden="true">1.1.1.</strong> DashAnyCompletion</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/dashCompletion.html"><strong aria-hidden="true">1.1.2.</strong> DashCompletion</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/flagCompletion.html"><strong aria-hidden="true">1.1.3.</strong> FlagCompletion</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/positionalAnyCompletion.html"><strong aria-hidden="true">1.1.4.</strong> PositionalAnyCompletion</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/positionalCompletion.html"><strong aria-hidden="true">1.1.5.</strong> PositionalCompletion</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/preInvoke.html"><strong aria-hidden="true">1.1.6.</strong> PreInvoke</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/preRun.html"><strong aria-hidden="true">1.1.7.</strong> PreRun</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/snippet.html"><strong aria-hidden="true">1.1.8.</strong> Snippet</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/gen/standalone.html"><strong aria-hidden="true">1.1.9.</strong> Standalone</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action.html"><strong aria-hidden="true">1.2.</strong> Action</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/cache.html"><strong aria-hidden="true">1.2.1.</strong> Cache</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/chdir.html"><strong aria-hidden="true">1.2.2.</strong> Chdir</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/chdirF.html"><strong aria-hidden="true">1.2.3.</strong> ChdirF</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/escape.html"><strong aria-hidden="true">1.2.4.</strong> Escape</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/filter.html"><strong aria-hidden="true">1.2.5.</strong> Filter</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/filterArgs.html"><strong aria-hidden="true">1.2.6.</strong> FilterArgs</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/filterParts.html"><strong aria-hidden="true">1.2.7.</strong> FilterParts</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/invoke.html"><strong aria-hidden="true">1.2.8.</strong> Invoke</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/list.html"><strong aria-hidden="true">1.2.9.</strong> List</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/multiParts.html"><strong aria-hidden="true">1.2.10.</strong> MultiParts</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/multiPartsP.html"><strong aria-hidden="true">1.2.11.</strong> MultiPartsP</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/noPrefix.html"><strong aria-hidden="true">1.2.12.</strong> NoPrefix</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/noSpace.html"><strong aria-hidden="true">1.2.13.</strong> NoSpace</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/prefix.html"><strong aria-hidden="true">1.2.14.</strong> Prefix</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/quote.html"><strong aria-hidden="true">1.2.15.</strong> Quote</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/retain.html"><strong aria-hidden="true">1.2.16.</strong> Retain</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/shift.html"><strong aria-hidden="true">1.2.17.</strong> Shift</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/split.html"><strong aria-hidden="true">1.2.18.</strong> Split</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/splitP.html"><strong aria-hidden="true">1.2.19.</strong> SplitP</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/style.html"><strong aria-hidden="true">1.2.20.</strong> Style</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/styleF.html"><strong aria-hidden="true">1.2.21.</strong> StyleF</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/styleR.html"><strong aria-hidden="true">1.2.22.</strong> StyleR</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/suffix.html"><strong aria-hidden="true">1.2.23.</strong> Suffix</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/suppress.html"><strong aria-hidden="true">1.2.24.</strong> Suppress</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/tag.html"><strong aria-hidden="true">1.2.25.</strong> Tag</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/tagF.html"><strong aria-hidden="true">1.2.26.</strong> TagF</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/timeout.html"><strong aria-hidden="true">1.2.27.</strong> Timeout</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/unique.html"><strong aria-hidden="true">1.2.28.</strong> Unique</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/uniqueList.html"><strong aria-hidden="true">1.2.29.</strong> UniqueList</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/uniqueListF.html"><strong aria-hidden="true">1.2.30.</strong> UniqueListF</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/unless.html"><strong aria-hidden="true">1.2.31.</strong> Unless</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/unlessF.html"><strong aria-hidden="true">1.2.32.</strong> UnlessF</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/usage.html"><strong aria-hidden="true">1.2.33.</strong> Usage</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/action/usageF.html"><strong aria-hidden="true">1.2.34.</strong> UsageF</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction.html"><strong aria-hidden="true">1.3.</strong> InvokedAction</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/escape.html"><strong aria-hidden="true">1.3.1.</strong> Escape</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/filter.html"><strong aria-hidden="true">1.3.2.</strong> Filter</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/merge.html"><strong aria-hidden="true">1.3.3.</strong> Merge</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/prefix.html"><strong aria-hidden="true">1.3.4.</strong> Prefix</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/quote.html"><strong aria-hidden="true">1.3.5.</strong> Quote</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/retain.html"><strong aria-hidden="true">1.3.6.</strong> Retain</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/suffix.html"><strong aria-hidden="true">1.3.7.</strong> Suffix</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/toA.html"><strong aria-hidden="true">1.3.8.</strong> ToA</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedAction/toMultiPartsA.html"><strong aria-hidden="true">1.3.9.</strong> ToMultiPartsA</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions.html"><strong aria-hidden="true">1.4.</strong> DefaultActions</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionCallback.html"><strong aria-hidden="true">1.4.1.</strong> ActionCallback</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionCobra.html"><strong aria-hidden="true">1.4.2.</strong> ActionCobra</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionCommands.html"><strong aria-hidden="true">1.4.3.</strong> ActionCommands</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionDirectories.html"><strong aria-hidden="true">1.4.4.</strong> ActionDirectories</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionExecCommand.html"><strong aria-hidden="true">1.4.5.</strong> ActionExecCommand</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionExecCommandE.html"><strong aria-hidden="true">1.4.6.</strong> ActionExecCommandE</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionExecutables.html"><strong aria-hidden="true">1.4.7.</strong> ActionExecutables</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionExecute.html"><strong aria-hidden="true">1.4.8.</strong> ActionExecute</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionFiles.html"><strong aria-hidden="true">1.4.9.</strong> ActionFiles</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionImport.html"><strong aria-hidden="true">1.4.10.</strong> ActionImport</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionMessage.html"><strong aria-hidden="true">1.4.11.</strong> ActionMessage</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionMultiParts.html"><strong aria-hidden="true">1.4.12.</strong> ActionMultiParts</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionMultiPartsN.html"><strong aria-hidden="true">1.4.13.</strong> ActionMultiPartsN</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionPositional.html"><strong aria-hidden="true">1.4.14.</strong> ActionPositional</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionStyleConfig.html"><strong aria-hidden="true">1.4.15.</strong> ActionStyleConfig</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionStyledValues.html"><strong aria-hidden="true">1.4.16.</strong> ActionStyledValues</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionStyledValuesDescribed.html"><strong aria-hidden="true">1.4.17.</strong> ActionStyledValuesDescribed</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionStyles.html"><strong aria-hidden="true">1.4.18.</strong> ActionStyles</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionValues.html"><strong aria-hidden="true">1.4.19.</strong> ActionValues</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/defaultActions/actionValuesDescribed.html"><strong aria-hidden="true">1.4.20.</strong> ActionValuesDescribed</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/customActions.html"><strong aria-hidden="true">1.5.</strong> CustomActions</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/context.html"><strong aria-hidden="true">1.6.</strong> Context</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/context/abs.html"><strong aria-hidden="true">1.6.1.</strong> Abs</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/context/command.html"><strong aria-hidden="true">1.6.2.</strong> Command</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/context/envsubst.html"><strong aria-hidden="true">1.6.3.</strong> Envsubst</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/context/getEnv.html"><strong aria-hidden="true">1.6.4.</strong> GetEnv</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/context/lookupEnv.html"><strong aria-hidden="true">1.6.5.</strong> LookupEnv</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/context/setEnv.html"><strong aria-hidden="true">1.6.6.</strong> SetEnv</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/batch.html"><strong aria-hidden="true">1.7.</strong> Batch</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/batch/invoke.html"><strong aria-hidden="true">1.7.1.</strong> Invoke</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/batch/ToA.html"><strong aria-hidden="true">1.7.2.</strong> ToA</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedBatch.html"><strong aria-hidden="true">1.8.</strong> InvokedBatch</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/invokedBatch/merge.html"><strong aria-hidden="true">1.8.1.</strong> Merge</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/export.html"><strong aria-hidden="true">1.9.</strong> Export</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/command.html"><strong aria-hidden="true">1.10.</strong> Command</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/command/group.html"><strong aria-hidden="true">1.10.1.</strong> Group</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/standalone.html"><strong aria-hidden="true">1.11.</strong> Standalone</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/standalone/carapace-parse.html"><strong aria-hidden="true">1.11.1.</strong> carapace-parse</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/standalone/pflag.html"><strong aria-hidden="true">1.11.2.</strong> pflag</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/sandbox.html"><strong aria-hidden="true">1.12.</strong> Sandbox</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/clearCache.html"><strong aria-hidden="true">1.12.1.</strong> ClearCache</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/env.html"><strong aria-hidden="true">1.12.2.</strong> Env</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/files.html"><strong aria-hidden="true">1.12.3.</strong> Files</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/keep.html"><strong aria-hidden="true">1.12.4.</strong> Keep</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/newContext.html"><strong aria-hidden="true">1.12.5.</strong> NewContext</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/reply.html"><strong aria-hidden="true">1.12.6.</strong> Reply</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/reply/with.html"><strong aria-hidden="true">1.12.6.1.</strong> With</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/run.html"><strong aria-hidden="true">1.12.7.</strong> Run</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/expect.html"><strong aria-hidden="true">1.12.7.1.</strong> Expect</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/expectNot.html"><strong aria-hidden="true">1.12.7.2.</strong> ExpectNot</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/output.html"><strong aria-hidden="true">1.12.7.3.</strong> Output</a></span></li></ol></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="carapace/multiCompleter.html"><strong aria-hidden="true">1.13.</strong> Multi Completer</a></span></li></ol><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development.html"><strong aria-hidden="true">2.</strong> development</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/testing.html"><strong aria-hidden="true">2.1.</strong> Testing</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/asciinema.html"><strong aria-hidden="true">2.2.</strong> Asciinema</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs.html"><strong aria-hidden="true">2.3.</strong> VHS</a><a class="chapter-fold-toggle"><div>❱</div></a></span><ol class="section"><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/word.html"><strong aria-hidden="true">2.3.1.</strong> Word</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/nospace.html"><strong aria-hidden="true">2.3.2.</strong> Nospace</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/doubleQuote.html"><strong aria-hidden="true">2.3.3.</strong> Double Quote</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/doubleQuoteOpen.html"><strong aria-hidden="true">2.3.4.</strong> Double Quote Open</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/doubleQuotePartial.html"><strong aria-hidden="true">2.3.5.</strong> Double Quote Partial</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/noPrefix.html"><strong aria-hidden="true">2.3.6.</strong> NoPrefix</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/singleQuote.html"><strong aria-hidden="true">2.3.7.</strong> Single Quote</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/singleQuoteOpen.html"><strong aria-hidden="true">2.3.8.</strong> Single Quote Open</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/singleQuotePartial.html"><strong aria-hidden="true">2.3.9.</strong> Single Quote Partial</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/escape.html"><strong aria-hidden="true">2.3.10.</strong> Escape</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/redirect.html"><strong aria-hidden="true">2.3.11.</strong> Redirect</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/special.html"><strong aria-hidden="true">2.3.12.</strong> Special</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/specialDoubleQuote.html"><strong aria-hidden="true">2.3.13.</strong> Special Double Quote</a></span></li><li class="chapter-item "><span class="chapter-link-wrapper"><a href="development/vhs/specialSingleQuote.html"><strong aria-hidden="true">2.3.14.</strong> Special Single Quote</a></span></li></ol></li></ol></li></ol>';
        // Set the current, active page, and reveal it if it's hidden
        let current_page = document.location.href.toString().split('#')[0].split('?')[0];
        if (current_page.endsWith('/')) {
            current_page += 'index.html';
        }
        const links = Array.prototype.slice.call(this.querySelectorAll('a'));
        const l = links.length;
        for (let i = 0; i < l; ++i) {
            const link = links[i];
            const href = link.getAttribute('href');
            if (href && !href.startsWith('#') && !/^(?:[a-z+]+:)?\/\//.test(href)) {
                link.href = path_to_root + href;
            }
            // The 'index' page is supposed to alias the first chapter in the book.
            // Check both with and without the '.html' suffix to be robust against pretty URLs
            if (link.href.replace(/\.html$/, '') === current_page.replace(/\.html$/, '')
                || i === 0
                && path_to_root === ''
                && current_page.endsWith('/index.html')) {
                link.classList.add('active');
                let parent = link.parentElement;
                while (parent) {
                    if (parent.tagName === 'LI' && parent.classList.contains('chapter-item')) {
                        parent.classList.add('expanded');
                    }
                    parent = parent.parentElement;
                }
            }
        }
        // Track and set sidebar scroll position
        this.addEventListener('click', e => {
            if (e.target.tagName === 'A') {
                const clientRect = e.target.getBoundingClientRect();
                const sidebarRect = this.getBoundingClientRect();
                sessionStorage.setItem('sidebar-scroll-offset', clientRect.top - sidebarRect.top);
            }
        }, { passive: true });
        const sidebarScrollOffset = sessionStorage.getItem('sidebar-scroll-offset');
        sessionStorage.removeItem('sidebar-scroll-offset');
        if (sidebarScrollOffset !== null) {
            // preserve sidebar scroll position when navigating via links within sidebar
            const activeSection = this.querySelector('.active');
            if (activeSection) {
                const clientRect = activeSection.getBoundingClientRect();
                const sidebarRect = this.getBoundingClientRect();
                const currentOffset = clientRect.top - sidebarRect.top;
                this.scrollTop += currentOffset - parseFloat(sidebarScrollOffset);
            }
        } else {
            // scroll sidebar to current active section when navigating via
            // 'next/previous chapter' buttons
            const activeSection = document.querySelector('#mdbook-sidebar .active');
            if (activeSection) {
                activeSection.scrollIntoView({ block: 'center' });
            }
        }
        // Toggle buttons
        const sidebarAnchorToggles = document.querySelectorAll('.chapter-fold-toggle');
        function toggleSection(ev) {
            ev.currentTarget.parentElement.parentElement.classList.toggle('expanded');
        }
        Array.from(sidebarAnchorToggles).forEach(el => {
            el.addEventListener('click', toggleSection);
        });
    }
}
window.customElements.define('mdbook-sidebar-scrollbox', MDBookSidebarScrollbox);


// ---------------------------------------------------------------------------
// Support for dynamically adding headers to the sidebar.

(function() {
    // This is used to detect which direction the page has scrolled since the
    // last scroll event.
    let lastKnownScrollPosition = 0;
    // This is the threshold in px from the top of the screen where it will
    // consider a header the "current" header when scrolling down.
    const defaultDownThreshold = 150;
    // Same as defaultDownThreshold, except when scrolling up.
    const defaultUpThreshold = 300;
    // The threshold is a virtual horizontal line on the screen where it
    // considers the "current" header to be above the line. The threshold is
    // modified dynamically to handle headers that are near the bottom of the
    // screen, and to slightly offset the behavior when scrolling up vs down.
    let threshold = defaultDownThreshold;
    // This is used to disable updates while scrolling. This is needed when
    // clicking the header in the sidebar, which triggers a scroll event. It
    // is somewhat finicky to detect when the scroll has finished, so this
    // uses a relatively dumb system of disabling scroll updates for a short
    // time after the click.
    let disableScroll = false;
    // Array of header elements on the page.
    let headers;
    // Array of li elements that are initially collapsed headers in the sidebar.
    // I'm not sure why eslint seems to have a false positive here.
    // eslint-disable-next-line prefer-const
    let headerToggles = [];
    // This is a debugging tool for the threshold which you can enable in the console.
    let thresholdDebug = false;

    // Updates the threshold based on the scroll position.
    function updateThreshold() {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const windowHeight = window.innerHeight;
        const documentHeight = document.documentElement.scrollHeight;

        // The number of pixels below the viewport, at most documentHeight.
        // This is used to push the threshold down to the bottom of the page
        // as the user scrolls towards the bottom.
        const pixelsBelow = Math.max(0, documentHeight - (scrollTop + windowHeight));
        // The number of pixels above the viewport, at least defaultDownThreshold.
        // Similar to pixelsBelow, this is used to push the threshold back towards
        // the top when reaching the top of the page.
        const pixelsAbove = Math.max(0, defaultDownThreshold - scrollTop);
        // How much the threshold should be offset once it gets close to the
        // bottom of the page.
        const bottomAdd = Math.max(0, windowHeight - pixelsBelow - defaultDownThreshold);
        let adjustedBottomAdd = bottomAdd;

        // Adjusts bottomAdd for a small document. The calculation above
        // assumes the document is at least twice the windowheight in size. If
        // it is less than that, then bottomAdd needs to be shrunk
        // proportional to the difference in size.
        if (documentHeight < windowHeight * 2) {
            const maxPixelsBelow = documentHeight - windowHeight;
            const t = 1 - pixelsBelow / Math.max(1, maxPixelsBelow);
            const clamp = Math.max(0, Math.min(1, t));
            adjustedBottomAdd *= clamp;
        }

        let scrollingDown = true;
        if (scrollTop < lastKnownScrollPosition) {
            scrollingDown = false;
        }

        if (scrollingDown) {
            // When scrolling down, move the threshold up towards the default
            // downwards threshold position. If near the bottom of the page,
            // adjustedBottomAdd will offset the threshold towards the bottom
            // of the page.
            const amountScrolledDown = scrollTop - lastKnownScrollPosition;
            const adjustedDefault = defaultDownThreshold + adjustedBottomAdd;
            threshold = Math.max(adjustedDefault, threshold - amountScrolledDown);
        } else {
            // When scrolling up, move the threshold down towards the default
            // upwards threshold position. If near the bottom of the page,
            // quickly transition the threshold back up where it normally
            // belongs.
            const amountScrolledUp = lastKnownScrollPosition - scrollTop;
            const adjustedDefault = defaultUpThreshold - pixelsAbove
                + Math.max(0, adjustedBottomAdd - defaultDownThreshold);
            threshold = Math.min(adjustedDefault, threshold + amountScrolledUp);
        }

        if (documentHeight <= windowHeight) {
            threshold = 0;
        }

        if (thresholdDebug) {
            const id = 'mdbook-threshold-debug-data';
            let data = document.getElementById(id);
            if (data === null) {
                data = document.createElement('div');
                data.id = id;
                data.style.cssText = `
                    position: fixed;
                    top: 50px;
                    right: 10px;
                    background-color: 0xeeeeee;
                    z-index: 9999;
                    pointer-events: none;
                `;
                document.body.appendChild(data);
            }
            data.innerHTML = `
                <table>
                  <tr><td>documentHeight</td><td>${documentHeight.toFixed(1)}</td></tr>
                  <tr><td>windowHeight</td><td>${windowHeight.toFixed(1)}</td></tr>
                  <tr><td>scrollTop</td><td>${scrollTop.toFixed(1)}</td></tr>
                  <tr><td>pixelsAbove</td><td>${pixelsAbove.toFixed(1)}</td></tr>
                  <tr><td>pixelsBelow</td><td>${pixelsBelow.toFixed(1)}</td></tr>
                  <tr><td>bottomAdd</td><td>${bottomAdd.toFixed(1)}</td></tr>
                  <tr><td>adjustedBottomAdd</td><td>${adjustedBottomAdd.toFixed(1)}</td></tr>
                  <tr><td>scrollingDown</td><td>${scrollingDown}</td></tr>
                  <tr><td>threshold</td><td>${threshold.toFixed(1)}</td></tr>
                </table>
            `;
            drawDebugLine();
        }

        lastKnownScrollPosition = scrollTop;
    }

    function drawDebugLine() {
        if (!document.body) {
            return;
        }
        const id = 'mdbook-threshold-debug-line';
        const existingLine = document.getElementById(id);
        if (existingLine) {
            existingLine.remove();
        }
        const line = document.createElement('div');
        line.id = id;
        line.style.cssText = `
            position: fixed;
            top: ${threshold}px;
            left: 0;
            width: 100vw;
            height: 2px;
            background-color: red;
            z-index: 9999;
            pointer-events: none;
        `;
        document.body.appendChild(line);
    }

    function mdbookEnableThresholdDebug() {
        thresholdDebug = true;
        updateThreshold();
        drawDebugLine();
    }

    window.mdbookEnableThresholdDebug = mdbookEnableThresholdDebug;

    // Updates which headers in the sidebar should be expanded. If the current
    // header is inside a collapsed group, then it, and all its parents should
    // be expanded.
    function updateHeaderExpanded(currentA) {
        // Add expanded to all header-item li ancestors.
        let current = currentA.parentElement;
        while (current) {
            if (current.tagName === 'LI' && current.classList.contains('header-item')) {
                current.classList.add('expanded');
            }
            current = current.parentElement;
        }
    }

    // Updates which header is marked as the "current" header in the sidebar.
    // This is done with a virtual Y threshold, where headers at or below
    // that line will be considered the current one.
    function updateCurrentHeader() {
        if (!headers || !headers.length) {
            return;
        }

        // Reset the classes, which will be rebuilt below.
        const els = document.getElementsByClassName('current-header');
        for (const el of els) {
            el.classList.remove('current-header');
        }
        for (const toggle of headerToggles) {
            toggle.classList.remove('expanded');
        }

        // Find the last header that is above the threshold.
        let lastHeader = null;
        for (const header of headers) {
            const rect = header.getBoundingClientRect();
            if (rect.top <= threshold) {
                lastHeader = header;
            } else {
                break;
            }
        }
        if (lastHeader === null) {
            lastHeader = headers[0];
            const rect = lastHeader.getBoundingClientRect();
            const windowHeight = window.innerHeight;
            if (rect.top >= windowHeight) {
                return;
            }
        }

        // Get the anchor in the summary.
        const href = '#' + lastHeader.id;
        const a = [...document.querySelectorAll('.header-in-summary')]
            .find(element => element.getAttribute('href') === href);
        if (!a) {
            return;
        }

        a.classList.add('current-header');

        updateHeaderExpanded(a);
    }

    // Updates which header is "current" based on the threshold line.
    function reloadCurrentHeader() {
        if (disableScroll) {
            return;
        }
        updateThreshold();
        updateCurrentHeader();
    }


    // When clicking on a header in the sidebar, this adjusts the threshold so
    // that it is located next to the header. This is so that header becomes
    // "current".
    function headerThresholdClick(event) {
        // See disableScroll description why this is done.
        disableScroll = true;
        setTimeout(() => {
            disableScroll = false;
        }, 100);
        // requestAnimationFrame is used to delay the update of the "current"
        // header until after the scroll is done, and the header is in the new
        // position.
        requestAnimationFrame(() => {
            requestAnimationFrame(() => {
                // Closest is needed because if it has child elements like <code>.
                const a = event.target.closest('a');
                const href = a.getAttribute('href');
                const targetId = href.substring(1);
                const targetElement = document.getElementById(targetId);
                if (targetElement) {
                    threshold = targetElement.getBoundingClientRect().bottom;
                    updateCurrentHeader();
                }
            });
        });
    }

    // Takes the nodes from the given head and copies them over to the
    // destination, along with some filtering.
    function filterHeader(source, dest) {
        const clone = source.cloneNode(true);
        clone.querySelectorAll('mark').forEach(mark => {
            mark.replaceWith(...mark.childNodes);
        });
        dest.append(...clone.childNodes);
    }

    // Scans page for headers and adds them to the sidebar.
    document.addEventListener('DOMContentLoaded', function() {
        const activeSection = document.querySelector('#mdbook-sidebar .active');
        if (activeSection === null) {
            return;
        }

        const main = document.getElementsByTagName('main')[0];
        headers = Array.from(main.querySelectorAll('h2, h3, h4, h5, h6'))
            .filter(h => h.id !== '' && h.children.length && h.children[0].tagName === 'A');

        if (headers.length === 0) {
            return;
        }

        // Build a tree of headers in the sidebar.

        const stack = [];

        const firstLevel = parseInt(headers[0].tagName.charAt(1));
        for (let i = 1; i < firstLevel; i++) {
            const ol = document.createElement('ol');
            ol.classList.add('section');
            if (stack.length > 0) {
                stack[stack.length - 1].ol.appendChild(ol);
            }
            stack.push({level: i + 1, ol: ol});
        }

        // The level where it will start folding deeply nested headers.
        const foldLevel = 3;

        for (let i = 0; i < headers.length; i++) {
            const header = headers[i];
            const level = parseInt(header.tagName.charAt(1));

            const currentLevel = stack[stack.length - 1].level;
            if (level > currentLevel) {
                // Begin nesting to this level.
                for (let nextLevel = currentLevel + 1; nextLevel <= level; nextLevel++) {
                    const ol = document.createElement('ol');
                    ol.classList.add('section');
                    const last = stack[stack.length - 1];
                    const lastChild = last.ol.lastChild;
                    // Handle the case where jumping more than one nesting
                    // level, which doesn't have a list item to place this new
                    // list inside of.
                    if (lastChild) {
                        lastChild.appendChild(ol);
                    } else {
                        last.ol.appendChild(ol);
                    }
                    stack.push({level: nextLevel, ol: ol});
                }
            } else if (level < currentLevel) {
                while (stack.length > 1 && stack[stack.length - 1].level > level) {
                    stack.pop();
                }
            }

            const li = document.createElement('li');
            li.classList.add('header-item');
            li.classList.add('expanded');
            if (level < foldLevel) {
                li.classList.add('expanded');
            }
            const span = document.createElement('span');
            span.classList.add('chapter-link-wrapper');
            const a = document.createElement('a');
            span.appendChild(a);
            a.href = '#' + header.id;
            a.classList.add('header-in-summary');
            filterHeader(header.children[0], a);
            a.addEventListener('click', headerThresholdClick);
            const nextHeader = headers[i + 1];
            if (nextHeader !== undefined) {
                const nextLevel = parseInt(nextHeader.tagName.charAt(1));
                if (nextLevel > level && level >= foldLevel) {
                    const toggle = document.createElement('a');
                    toggle.classList.add('chapter-fold-toggle');
                    toggle.classList.add('header-toggle');
                    toggle.addEventListener('click', () => {
                        li.classList.toggle('expanded');
                    });
                    const toggleDiv = document.createElement('div');
                    toggleDiv.textContent = '❱';
                    toggle.appendChild(toggleDiv);
                    span.appendChild(toggle);
                    headerToggles.push(li);
                }
            }
            li.appendChild(span);

            const currentParent = stack[stack.length - 1];
            currentParent.ol.appendChild(li);
        }

        const onThisPage = document.createElement('div');
        onThisPage.classList.add('on-this-page');
        onThisPage.append(stack[0].ol);
        const activeItemSpan = activeSection.parentElement;
        activeItemSpan.after(onThisPage);
    });

    document.addEventListener('DOMContentLoaded', reloadCurrentHeader);
    document.addEventListener('scroll', reloadCurrentHeader, { passive: true });
})();

