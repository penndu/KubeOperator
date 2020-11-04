import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Manifest} from '../manifest';
import {ManifestService} from '../manifest.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-manifest-list',
    templateUrl: './manifest-list.component.html',
    styleUrls: ['./manifest-list.component.css']
})
export class ManifestListComponent implements OnInit {
    constructor(private manifestService: ManifestService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    items: Manifest[] = [];
    loading = false;
    @Output() detailEvent = new EventEmitter<Manifest>();

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.manifestService.list().subscribe(data => {
            this.items = data;
        });
    }

    onDetail(item: Manifest) {
        console.log(item);
        this.detailEvent.emit(item);
    }

    update(item: Manifest) {
        console.log(item);
        const updateItem = new Manifest();
        Object.assign(updateItem, item);
        updateItem.isActive = !item.isActive;
        this.manifestService.update(updateItem).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
