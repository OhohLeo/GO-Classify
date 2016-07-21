var ws;
var h2_collection_name;
var nav_collections;
var ul_collections;

var menu = {
    name: "",
    onSelect: function(name) {
        if (this.name == name) {
            return;
        }
        $('div.menu').addClass('hide');
        $('div#'+name).removeClass('hide');
        this.name = name;
    },
};

function Collection (name) {
    this.name = name;
};

Collection.prototype.onRename = function(name) {
    this.name = name;
};

var collections = {
    current: null,
    references: null,
    list: {},
    getReferences: function(cb) {

        if (this.references != null) {
            if (cb) {
                cb(this.references);
            }
            return
        }

        send('GET', 'references', undefined,
             function(body, status) {
                 if (cb) {
                     cb(body)
                 }
                 this.references = body;
             });
    },
    setTypes: function(types) {
        for (var i in types) {
            $('select#collection-create-type').append(
                $('<option>').val(i)
                             .text(types[i]));
        }
    },
    setWebsites: function(websites) {
        for (var i in websites) {
            var website = websites[i];
            $('ul#collection-create-websites').append(
                $('<label>').text(website),
                $('<input>').attr('type', 'checkbox',
                                  'name', website));
            $('ul#collection-modify-websites').append(
                $('<label>').text(website),
                $('<input>').attr('type', 'checkbox',
                                  'name', website));
        }
    },
    onSelect: function(name) {

        collection = this.list[name];
        if (collection == null) {
            alert("Unknown collection!");
            return;
        }

        if (this.current != null && this.current.name == name) {
            return;
        }

        $('h2.with-collection').removeClass('hide');
        h2_collection_name.text(name);
        ul_collections.children().removeClass('selected');
        $('ul#collections li:contains("' + name + '")').addClass('selected');
        nav_collections.slideToggle();
        menu.onSelect('collection');
        this.current = collection;
    },
    create: function(name, type, websites) {

        if (this.list[name] != null) {
            alert("Already existing collection!");
            return;
        }

        send('POST', 'collections', {
            name: name,
            type: type,
            websites: websites,
        }, function() {
            this.onCreate(name, type);
        });
    },
    onCreate: function(name, type) {

        if (this.list[name] != null) {
            alert("Already existing collection!");
            return;
        }

        this.list[name] = new Collection(name);

        ul_collections.append($('<li>').addClass("menu")
                                       .text(name));
        ul_collections.children().on('click', function() {
            collections.onSelect($(this).text());
        });

        nav_collections.slideDown();
    },
    onModify: function(name) {

        if (this.current == null) {
            alert("Not selected collection!");
            return;
        }

        if (this.current.name == name) {
            alert("Similar collection name!");
            return;
        }

        collection.onRename(name);
        h2_collection_name.text(name);
        $('ul#collections li.selected').text(name);
        nav_collections.slideDown();
    },
    onDelete: function(name) {

        if (this.current == null) {
            alert("Not selected collection!");
            return;
        }

        delete this.list[this.current.name];

        $('h2.with-collection').addClass('hide')
             .next().addClass('hide');
        h2_collection_name.text("");
        $('ul#collections li.selected').remove();
        nav_collections.slideDown();
    },
};

$("document").ready(function(){
    h2_collection_name = $('h2#collection-name');
    nav_collections = $('nav#collections');
    ul_collections = $('ul#collections');

    // get the collection types & websites available
    collections.getReferences(function (body) {
        collections.setTypes(body.types);
        collections.setWebsites(body.websites);
    });

    // get the collections
    collections.onCreate("Movies");
    collections.onCreate("Musics");

	//ws = start('ws://localhost:8080/ws');

    // click & display the main menus: Import, Export & Config
    $('li.menu').on('click', function() {
        menu.onSelect($(this).text().toLowerCase());
    });

    // click on the main collection : display collections
    h2_collection_name.on('click', function() {
        nav_collections.slideToggle();
    });

    // close & display the configuration title
    $('h2.config').on('click', function () {
        $(this).next().toggleClass('hide');
    });

    // create collection
    $('form#collection-create').submit(function (e) {
        e.preventDefault();

        var name = $('input#collection-create-name').val();
        if (name == '') {
            return;
        }

        var type = $('input#collection-create-type').val();
        var websites = $('input#collection-create-websites').val();

        collections.create(name, type, websites);
    });

    // modify collection
    $('form#collection-modify').submit(function (e) {
        e.preventDefault();

        var name = $('input#collection-modify-name').val();
        if (name == '' || name == this.name) {
            return;
        }

        var websites = $('input#collection-modify-websites').val();

        console.log(name);
        collections.onModify(name, websites);
    });

    $('form#collection-delete').submit(function (e) {
        e.preventDefault();
        collections.onDelete();
    });

	$('input#new-directory').on('change', function(){
		wsSend(ws, "new-directory", {
			path: $('input#newDirectory').val()
		});
	});
});

var onEvent = {
	newFile: function(data) {
		$("input#new-directory").after($("<p>").text(data.FullPath));
	},
}

function start(url)
{
	ws = new WebSocket(url);

	ws.onopen = function(e) {
		console.log('Connection OK');
	};

	var onClose = {
		1000: "Close OK",
		1001: "Server/Client DOWN",
		1002: "Protocol error",
		1003: "Type of data not accepted",
		1004: "Reserved",
		1005: "No status code",
		1006: "Connection was closed abnormally",
		1007: "Message not consistent",
		1008: "Message violates policy",
		1009: "Message too big",
		1010: "Negotiation failed",
		1011: "Unexpected condition",
		1015: "TLS handshake failed"
	};

	ws.onclose = function(e) {
		console.log('Connection close: %s', onClose[e.code]);
		setTimeout(function(){start(url)}, 5000);
	}

	ws.onmessage = function(e) {
		var obj = JSON.parse(e.data);
		console.log(' --> %o', obj);
		onEvent[obj.Type](obj.Data);
	};

	ws.onerror = function(e) {
		console.log(' !!! %o', e);
	};

	return ws
}

function wsSend(ws, type, data) {
	var json = JSON.stringify({type: type, data: data});
	console.log(' <-- %s', json);
	ws.send(json);
}

function send(method, path, body, success) {
    $.ajax({
        method: method,
        url: 'http://localhost:8080/'+path,
        context: body,
        success: success,
        error: function() {

        },
        dataType: 'json',
    });
}
