function filterViewResults() {
    $.ajax({
        method: "GET",
        url: "/books/" + $("#filter-view-results").serialize().substring(7),
        success: rebuildBookCollection
    })
}

function sortBooks(column) {
    $.ajax({
        method: "GET",
        url: "/books/" + column,
        success: rebuildBookCollection
    })
}

function rebuildBookCollection(result) {
    var books = JSON.parse(result);
    if (!books) return;

    $("#view-results").empty();

    books.forEach(appendBook)
}

function appendBook(book) {
    $("#view-results").append("<tr id='book-row-" + book.PK + "'><td>" + book.Title + "</td><td>" + book.Author + "</td><td>" + book.Classification +
            "</td><td><button class='delete-btn' onclick='deleteBook(" + book.PK + ")'>Delete</button></td></tr>");
}

function showSearchPage() {
    $("#search-page").show();
    $("#view-page").hide();
}

function showViewPage() {
    $("#view-page").show();
    $("#search-page").hide();
}

function deleteBook(pk) {
    $.ajax({
        url: "/books/" + pk,
        method: "DELETE",
        success: function() {
            $("#book-row-" + pk).remove();
        }
    })
}

function submitSearch() {
    $.ajax({
        url: "/search",
        method: "POST",
        data: $("#search-form").serialize(),
        success: function(rawData) {
            var parsed = JSON.parse(rawData);
            if (!parsed) return;

            var searchResults = $("#search-results");
            searchResults.empty();
            parsed.forEach(function(result) {
                var row = $("<tr><td>" + result.Title + "</td><td>" + result.Author + "</td><td>" + result.Year + "</td><td>" + result.ID + "</td></tr>");
                searchResults.append(row);
                row.on("click", function() {
                    $.ajax({
                        url: "/books/" + result.ID,
                        method: "PUT",
                        success: function(data) {
                            var book = JSON.parse(data);
                            if (!book) return;
                            appendBook(book);
                        }
                    })
                })
            })
        }
    });
    return false;
}
