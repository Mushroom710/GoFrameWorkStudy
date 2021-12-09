window.onload = function () {
    //使用事件代理。处理点击某个分类时，展示对应的分类的文章信息
    $('#select a').on('click', function () {
        // alert($(this).text());
        // attr可以对标签的属性进行操作
        // href的跳转地址需要写全
        $(this).attr('href', 'http://localhost:8080/?select=' + $(this).text())
    })
    $('#pic').click(function () {
        $.ajax({
            url: 'http://localhost:8080/ShowPicture',
            type: 'get',
            success: function (res) {
                const oRoot = document.getElementById('root')
                oRoot.innerHTML = res
            }
        })
    })
    $('#about').click(function () {
        $.ajax({
            url: 'http://localhost:8080/ShowAbout',
            type: 'get',
            success: function (res) {
                const oRoot = document.getElementById('root')
                oRoot.innerHTML = res
            }
        })
    })
    $('#root .articleAjax').on('click', function () {
        var num = $(this).attr('id')
        $.ajax({
            url: 'http://localhost:8080/PageArticle/' + num,
            type: 'get',
            success: function (res) {
                const oRoot = document.getElementById('root')
                oRoot.innerHTML = res
            }
        })
    })
    $('.star a').on('click', function () {
        var num = $(this).attr('id')
        $.ajax({
            url: 'http://localhost:8080/PageArticle/' + num,
            type: 'get',
            success: function (res) {
                const oRoot = document.getElementById('root')
                oRoot.innerHTML = res
            }
        })
    })
}