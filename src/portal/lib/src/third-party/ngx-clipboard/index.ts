import { ClipboardDirective } from './clipboard.directive';
import { CLIPBOARD_SERVICE_PROVIDER } from './clipboard.service';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { WindowTokenModule } from '../ngx-window-token/index';
export * from './clipboard.directive';
export * from './clipboard.service';

@NgModule({
    imports: [CommonModule, WindowTokenModule],
    declarations: [ClipboardDirective],
    exports: [ClipboardDirective],
    providers: [CLIPBOARD_SERVICE_PROVIDER]
})
export class ClipboardModule { }
