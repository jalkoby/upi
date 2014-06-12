document.addEventListener("DOMContentLoaded", function() {
  (function() {
    var links = document.querySelectorAll("a[data-token]"), length = links.length, i;
    for(i = 0; i < length; i++) {
      links[i].addEventListener("click", function(e) {
        e.preventDefault();
        document.getElementById("preview").hidden = false;
        document.getElementById("link").innerText = e.target.href;
      });
    }
  })();

  (function() {
    var links = document.querySelectorAll("a[data-method]"), length = links.length, i;
    for(i = 0; i < length; i++) {
      (function(link) {
        link.addEventListener("click", function(e) {
          e.preventDefault();
          var xhr = new XMLHttpRequest();
          xhr.onreadystatechange = function() {
            if (xhr.readyState == 4 ) {
              if(xhr.status == 204) { location.href = xhr.getResponseHeader("Location"); }
            }
          }
          xhr.open(link.dataset["method"], link.href, true);
          xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
          xhr.send();
        });
      })(links[i]);
    }
  })();
});
