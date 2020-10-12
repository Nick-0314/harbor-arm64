## Developer's Guide for Internationalization (i18n)

*NOTE: All the files you created should use UTF-8 encoding.*

### Steps to localize the UI in your language

1. In the folder `src/portal/src/i18n/lang`, copy json file `en-us-lang.json` to a new file and rename it to `<language>-<locale>-lang.json` .

    The file contains a JSON object including all the key-value pairs of UI strings:
    ```
    {
        "APP_TITLE": {
            "VMW_HARBOR": "Harbor",
            "HARBOR": "Harbor",
            ...
        },
        ...
    }
    ```  
    In the file `<language>-<locale>-lang.json`, translate all the values into your language. Do not change any keys.

2. After creating your language file, you should add it to the language supporting list.

    Locate the file `src/portal/src/app/shared/shared.const.ts`.
    Append `<language>-<locale>` to the language supporting list:
    ```
    export const supportedLangs = ['en-us', 'zh-cn', '<language>-<locale>'];
    ```
    Define the language display name and append it to the name list:
    ```
    export const languageNames = {
        "en-us": "English",
        "zh-cn": "中文简体",
        "<language>-<locale>": "<DISPLAY_NAME>"
    };
    ```

    **NOTE: Don't miss the comma before the new key-value item you've added.**

3. Enable the new language in the view.

    Locate the file `src/portal/src/app/base/navigator/navigator.component.html` and then find the following code piece:
    ```
    <div class="dropdown-menu">
        <a href="javascript:void(0)" clrDropdownItem (click)='switchLanguage("en-us")' [class.lang-selected]='matchLang("en-us")'>English</a>
        <a href="javascript:void(0)" clrDropdownItem (click)='switchLanguage("zh-cn")' [class.lang-selected]='matchLang("zh-cn")'>中文简体</a>
    </div>
    ```
    Add new menu item for your language:
    ```
    <div class="dropdown-menu">
        <a href="javascript:void(0)" clrDropdownItem (click)='switchLanguage("en-us")' [class.lang-selected]='matchLang("en-us")'>English</a>
        <a href="javascript:void(0)" clrDropdownItem (click)='switchLanguage("zh-cn")' [class.lang-selected]='matchLang("zh-cn")'>中文简体</a>
        <a href="javascript:void(0)" clrDropdownItem (click)='switchLanguage("<language>-<locale>")' [class.lang-selected]='matchLang("<language>-<locale>")'>DISPLAY_NAME</a>
    </div>
    ```

4. Next, please refer [compile guideline](compile_guide.md) to rebuild and restart Harbor.
