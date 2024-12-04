import Graph from "./graph.js";

const graph = new Graph(document.getElementById('d3'));

//	DOM nodes containing important information
const PAGES_QUERY = 'table#pages > tbody > tr';
const TAGS_QUERY = 'table#tags > tbody > tr';

//	For fast lookup
const pages = new Map();	//	permalink is used as key
const tags = new Map();		//	title is used as key

//	fully massaged graph data
const data = {
	vertices: [],
	edges: []
};

const breadcrumb = document.querySelector('nav.breadcrumb');

//	event-handlers
graph.onMouseOver = (_, d) => {
	if (d.type === "tag") {
		breadcrumb.querySelectorAll('a.dynamic').forEach(el => el.remove());
		const newAnchor = document.createElement('a');
		newAnchor.setAttribute('href', d.permalink);
		newAnchor.innerText = d.title;
		newAnchor.classList.add('dynamic');
		breadcrumb.append(newAnchor);
	}
};
graph.onMouseOut = (ev, d) => {
	breadcrumb.querySelectorAll('a.dynamic').forEach(el => el.remove());
};
graph.onClick = (_, d) => {
	console.log(_,d);
	location.href = d.permalink;
};


const ingest = async () => {

	let tagsRows = document.querySelectorAll(TAGS_QUERY);
	let pagesRows = document.querySelectorAll(PAGES_QUERY);	

	while (tagsRows.length) {
		let thisRow = tagsRows.item(0);
		let cells = thisRow.querySelectorAll('td');
		let permalink = cells.item(0).innerText;
		let title = cells.item(1).innerText;
		let centrality = Number(cells.item(2).innerText);
		let type = "tag";
		let thisVertex = graph.toVertex({
			type,
			title,
			permalink,
			centrality
		});
		tags.set(title, thisVertex);
		data.vertices.push(thisVertex);
		thisRow.remove();
		tagsRows = document.querySelectorAll(TAGS_QUERY);
	}

	while (pagesRows.length) {
		let thisRow = pagesRows.item(0);
		let cells = thisRow.querySelectorAll('td');
		let permalink = cells.item(0).innerText;
		let title = cells.item(1).innerText;
		let section = cells.item(2).innerText;
		let tagsForThisPage = Array.from( cells.item(3).querySelectorAll('li') ).map(el => el.innerText);
		let type = "article";
		let centrality = 1;
		let thisVertex = graph.toVertex({
			type,
			section,
			title,
			permalink,
			centrality
		});
		pages.set(permalink, thisVertex);
		data.vertices.push(thisVertex);
		let edgesForThisPage = tagsForThisPage.map(word => {
			let type = "articleIsTaggedWith"
			let source = tags.get(word);
			let target = thisVertex;
			return graph.toEdge({source,target,type});
		});
		data.edges.push(...edgesForThisPage);
		thisRow.remove();
		pagesRows = document.querySelectorAll(PAGES_QUERY);
	}

	while (data.vertices.length) {
		await graph.addVertex(data.vertices.pop());
	}

	while (data.edges.length) {
		await graph.addEdge(data.edges.pop());
	}

	return graph;

};

ingest().then(graph => {
	const s = graph.simulation;
	s.alphaMin(0.001);
	s.alphaDecay(0.0228);
	s.force("r", graph.d3.forceRadial(d => {
		let radius = (Math.min(graph.width, graph.height) / 2) - 25;
		return Math.max(50,(radius - (d.centrality*50)));
	}, graph.width/2, graph.height/2));
	s.alpha(1).restart();
    console.log(data);
});


/*
fetch("miserables.json").then(response => response.json()).then(data => {

	//	massage vertices and edges for rubost uniqueness
	const vertices = data.nodes.map(node => {
		return graph.toVertex(node);
	});
	const edges = data.links.map(link => {
		const _id = uid();
		const source = vertices.filter(vertex => {
			return (link.source === vertex.id)
		}).pop();
		const target = vertices.filter(vertex => {
			return (link.target === vertex.id)
		}).pop();
		return {_id, source, target};
	});
	const graphData = {vertices, edges};

	const popAll = async () => {

		while (graphData.vertices.length) {
			let miser = graphData.vertices.pop();
			await graph.addVertex(miser);
		}

		while (graphData.edges.length) {
			let edge = graphData.edges.pop();
			await graph.addEdge(edge);
		}

		return graph;

	};

	popAll().then(graph => {
		const s = graph.simulation;
		s.alphaMin(0.001);
		s.alphaDecay(0.0228);
		s.restart();
	}).catch(console.error);

});
*/