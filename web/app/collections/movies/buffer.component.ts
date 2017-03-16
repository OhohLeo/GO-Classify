import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'
import { Movie } from './movie'

@Component({
    selector: 'buffer-movie',
    templateUrl: './buffer.component.html'
})

export class BufferMovieComponent implements OnInit, OnDestroy {

    @Input() movie: Movie

    ngOnInit() {
    }

    ngOnDestroy() {
    }
}
