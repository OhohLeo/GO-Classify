import { Component, NgZone, OnInit, ViewChild } from '@angular/core'

import { ApiService, CollectionStatus, Event } from './api.service'
import { ImportsService } from './imports/imports.service'
import { BufferService } from './buffer/buffer.service'
import { CollectionsService } from './collections/collections.service'

import { Collection } from './collections/collection'

import { ListCollectionsComponent } from './collections/list.component'
import { BufferComponent } from './buffer/buffer.component'
import { BufferItemComponent } from './buffer/item.component'
import { BufferItem } from './buffer/item'
import { Item } from './items/item'

declare var jQuery: any;

enum AppStatus {
    NONE = 0,
    COLLECTION,
    IMPORTS,
    EXPORTS,
    CONFIGS,
    BUFFER_ITEM,
}

@Component({
    selector: 'app',
    templateUrl: './app.component.html',
})

export class AppComponent implements OnInit {
    @ViewChild(ListCollectionsComponent) collections: ListCollectionsComponent
    @ViewChild(BufferComponent) buffer: BufferComponent

    public appStatus = AppStatus
    public status = AppStatus.NONE

    public title = "Classify"

    public collection: Collection
    private otherCollections: Collection[] = []

    public modalTitle: string
    public modalMsg: string

    private importsLoop: any
    private importsRunningNb: number

    private searchEnabled: boolean
    private filterEnabled: boolean

    private menuActive: boolean
    private bufferActive: boolean
    private searchActive: boolean
    private filterActive: boolean    
    private addCollectionActive: boolean = true
    
    @ViewChild(BufferItemComponent) bufferItemComponent: BufferItemComponent
    public bufferItem: BufferItem

    constructor(private zone: NgZone,
        private apiService: ApiService,
        private importsService: ImportsService,
        private bufferService: BufferService,
        private collectionsService: CollectionsService) { }

    ngOnInit() {

        // Import loop
        this.importsLoop = jQuery("i#imports-loop")
        this.importsRunningNb = 0;
	
        // Flux inscription
        this.apiService.getStream()
            .subscribe((e: Event) => {

                console.log("EVENT!", e)

                // Detect restart application
                if (e.event === "start") {
                    window.location.replace("/");
                    return;
                }

                // Import data
                if (new RegExp("^import").test(e.event)) {
                    this.handleImport(e);
                    return;
                }

                // Collections reception
                if (new RegExp("^collection").test(e.event)) {
                    this.handleCollection(e);
                    return;
                }
            })

        this.apiService.subscribeCollectionChange(
            (collection: Collection, status: CollectionStatus) => {

		console.log("SUBSCRIBE API", collection, status)
		
		switch (status) {
		case CollectionStatus.CREATED:
		case CollectionStatus.DELETED:
		    this.collections.refresh(false)
		    break
		}		
		
                if (status == CollectionStatus.DELETED
		    || collection === undefined) {
		    console.log("ON COLLECTION CHANGE")
		    this.onChangeCollection()
                    return
                }

		console.log("COLLECTION HANDLE")
		this.onCollection(collection)
            })
    }

    onReset() {
    
	// Reset display first
        this.onNewState(AppStatus.NONE)
	
	// Disable collection menu
	this.enableMenu(false)
	
	// Reset title name
	this.setTitle("")

	// Add collection icon
	this.zone.run(() => {
	    this.addCollectionActive = true
	})
    }
    
    onCollection(collection: Collection) {

	console.log("ON COLLECTION:", collection)
	
        this.title = collection.name
        this.collection = collection

	// Ignore current collection
	let otherCollections = this.collections.getCollections(collection)
	this.addCollectionActive = (otherCollections.length == 0)
	this.otherCollections = otherCollections
	    
	// Activate collection menu
	this.enableMenu(true)

	// Select collection nav
	this.selectNav("collection")
	
        this.onNewState(AppStatus.COLLECTION)
    }

    onImports() {
        this.onNewState(AppStatus.IMPORTS)
    }

    onExports() {
        this.onNewState(AppStatus.EXPORTS)
    }

    onConfigs() {
        this.onNewState(AppStatus.CONFIGS)
    }
    
    onBufferItem(bufferItem: BufferItem) {
        this.zone.run(() => {
            this.bufferItem = bufferItem
            this.onNewState(AppStatus.BUFFER_ITEM)
        })
    }

    onNewState(nextStatus: AppStatus) {

	if (nextStatus == AppStatus.COLLECTION) {
	    this.enableFilterAndSearch(true)
	} else {
	    this.enableFilterAndSearch(false)
	}
	
	// Enable nav indicator
	jQuery("li.indicator").css("height", "2px")
	
	this.resetCollectionState()
        this.status = nextStatus
    }
    
    onBuffer() {
        this.buffer.start()
    }

    setTitle(name: string) {

	if (name == "") {
	    name = "Classify"
	}
	
	this.zone.run(() => {
	    this.title = name
	})
    }

    enableMenu(status: boolean) {

	if (this.menuActive == status)
	    return
	
	this.zone.run(() => {
	    this.menuActive = status
	})
    }
    
    enableFilter(status: boolean) {
	this.zone.run(() => {
            this.filterEnabled = status
	    if (status == false) {
		this.onFilterClose()
	    }
	})
    }
    
    onFilter() {

	// Toggle filtering
	this.changeFilterState(!this.filterActive)
    }

    onFilterClose() {
	this.changeFilterState(false)
    }

    changeFilterState(newState: boolean) {
	this.zone.run(() => {
            this.filterActive = newState
	    this.class2toggle("li#filter", "active", newState)
	})
    }
    
    enableSearch(status: boolean) {
	
	this.zone.run(() => {

	    this.searchEnabled = status

	    if (status == false) {
		this.onSearchClose()
	    }
        })
    }

    onSearch() {
	
	// Toggle search
	this.changeSearchState(!this.searchActive)
    }
    
    onSearchClose() {
	this.changeSearchState(false)
    }
    
    changeSearchState(newState: boolean) {
	this.zone.run(() => {
            this.searchActive = newState
	    this.class2toggle("li#search", "active", newState)
	})
    }

    class2toggle(item: string, className: string, toggle: boolean) {

	if (toggle) {
	    jQuery(item).addClass(className)
	} else {
	    jQuery(item).removeClass(className)
	}
    }

    enableFilterAndSearch(status: boolean) {
	this.enableFilter(status)
	this.enableSearch(status)
    }
    
    handleImport(e: Event) {

        console.log("IMPORT?", e)

        // Send notifications to the imports list
        this.importsService.addEvent(e);

        // Display imports status
        if (new RegExp('status$').test(e.event)) {

            console.log("Status??")

            // Status 'TRUE': rotate refresh logo
            if (e.data) {
                this.importsLoop.addClass("rotation")
                this.importsRunningNb++
            }
            // Status 'FALSE'
            else if (this.importsRunningNb > 0) {
                this.importsRunningNb--
            }

            // No more imports : stop logo rotation
            if (this.importsRunningNb < 1) {
                this.importsLoop.removeClass("rotation")
            }
        }
    }

    handleCollection(e: Event) {

        let names = e.event.split("/");
        let size = names.length;

        if (size <= 1) {
            console.error("Invalid collection event '" + e.event + "'")
            return;
        }

        let collection = this.collectionsService.getCollection(names[1])

        let destination: string
        if (size > 2) {
            destination = names[2];
        }

        // Send notifications to the imports list
        switch (destination) {
            case "buffers":

            let bufferItem = new BufferItem(e.data)

            if (this.bufferItem != undefined
                && this.bufferItem.id == bufferItem.id) {

                this.bufferItemComponent.onUpdate(bufferItem)
                this.bufferItem = bufferItem
            }
            // collection.addBufferItem(new BufferItem(e.data))
	    
            break;

        case "items":
            collection.addItem(e.data)
            break;

	default:
	    console.error("Unhandled collection destination '" + destination + "'")
        }
    }

    onError(title: string, msg: string) {

        console.error(title + "error :" + msg)

        this.modalTitle = title + " error!"
        this.modalMsg = msg
        jQuery('#modal').openModal()
    }

    stopModal() {
        jQuery('#modal').closeModal()
    }

    selectNav(name: string) {
	jQuery("li#"+name).children().click()
    }
    
    disableNav() {

	// Unbold nav title
	jQuery("li.tab").children().removeClass("active")

	// Disable nav indicator
	jQuery("li.indicator").css("height", "0px")	    
    }
    
    // Display new collection
    onNewCollection() {

	this.onReset()
	
	if (this.collections) {
	    this.collections.onNewCollection()
        }
    }

    // Display collections list
    onChangeCollection() {

	console.log("CHANGE COLLECTION??")
	
	// Get collection nb
	let collectionNb = this.collections.nb()

	// If no collection exist : create one
	if (collectionNb == 0) {
	    this.onReset()
	    this.onNewCollection()
	    return
	}
	
	// If collections list exist
	if (collectionNb > 1) {
	    this.onReset()
	    this.collections.onChooseCollection(undefined)
	    return
	}

	// Otherwise select current collection
	if (this.collection) {
	    this.onCollection(this.collection)
	}
    }

    // Affiche la liste des collections à sélectionner si aucune
    // collection n'est actuellement sélectionnée
    resetCollectionState() {
        if (this.collections)
            this.collections.resetCollectionState()
    }
}
