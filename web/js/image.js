class LiveImage {
    constructor(containerId) {  
        this.m_containerId = containerId;
        this.m_container = document.getElementById(containerId);
    }

    toString() {
        this.m_container.innerHTML = "LiveImage";
    }
}

function ShowLiveImage(containerId) {
    this.liveImage = new LiveImage(containerId);

    this.liveImage.toString();
}