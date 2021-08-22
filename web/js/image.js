class LiveImage {
    constructor(containerId, source) {
        this.m_containerId = containerId;
        this.m_container = document.getElementById(containerId);
        this.m_source = source;
    }

    display() {
        var image = document.getElementById("liveimage");
        if (image == null) {
            image = document.createElement("img");
            image.id = "liveimage";
            image.src = "live";
            image.style.maxWidth = "100%";
            image.style.maxHeight = "100%";
            image.style.margin = "auto";
            this.m_container.appendChild(image);
        }
    }

    toString() {
        this.m_container.innerHTML = "LiveImage";
    }
}

function ShowLiveImage(containerId) {
    this.liveImage = new LiveImage(containerId);

    this.liveImage.display();
}