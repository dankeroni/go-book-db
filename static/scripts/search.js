function showSearchPage() {
    $("#search-page").show();
    $("#view-page").hide();
}

function showViewPage() {
    $("#view-page").show();
    $("#search-page").hide();
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
                        url: "/books/add?id=" + result.ID,
                        method: "GET",
                        success: function(data) {
                            var book = JSON.parse(data);
                            if (!book) return;
                            console.log(book.Classification.MostPopular);
                            $("#view-results").append("<tr><td>" + book.Title + "</td><td>" + book.Author + "</td><td>" + book.Classification.MostPopular + "</td></tr>");
                        }
                    })
                })
            })
        }
    });
    return false;
}
