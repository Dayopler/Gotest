# Doc pour dossier de formation

## Plan

Sur ces parties il est possible de broder en expliquant de facon large pour que ce soit le plus comprehensible possible
### Introduire sujet (explication large sans (trop) de technique pour que ce soit compréhensible par tous)
- expliquer d'ou vient le besoin
- expliquer la solution existante (vielle version non generique et fonctionnant avec une techno metier particuliere, les fichiers protobuf format sesame)
- expliquer (en gros et brodant) le principe du projet et pour qui ca va servir et pourquoi (les datascientist)

### Commencer a introiduire le fonctionne du principe et les techno utilisé
- Cesium (dire ce que c'est, par qui s'est utilisé et dans quel but on s'en sert)
- Parler d'openSkyNetwork pour la recolte de donnée et le besoin potentiel de vision d'information ou de probleme qui peuve etre resolu grace a l'appli
- React/typescript (dire que c'est deja un peu en interne pour react et les points fort et pourquoi ca, typescript car passage a typescript pour le front car js utilisé avant)


---
# Cette partie la est un peu (beaucoup ) plus technique et vient expliquer le front et le deploiement
## Principe
Seed-offline (systeme evaluation et enregistrement de données) est une applciation utilisé afin d'analyser un grand volume de données aeriennes. Le principe, dans un premier temps, est de Drag and Drop des fichiers de données en GeoJSON dans l'application afin d'obtenir d'unité aerienne sur un affichage sur une cartographie 3D de la planete Terre.

### Generique
Une conception générique des fichiers de données GeoJSON a été réalisé afin qu'un grand nombre de personne puisse, facilement, utiliser cette application.

Le but d'avoir fait une conception générique des fichiers ainsi que l'application est qu'on vient a utiliser des filtres sur ces données. Nous avons la possibilité de créer ses propres filtres via un espace de creation directement sur l'IHM. Il est egalement possible de sauvegarder ces filtres, l'application vient creer et telecharger un fichier JSON qui peut etre a nouveau dropper dans l'application afin obtenir directement ses filtres.
Il est possible de supprimer un filtre directement depuis l'application
Le coté générique est utile pour pouvoir filtrer sur n'importe quel champ, si l'utilisateur veut filtrer sur le champ `nombre reacteurs : 4` il faut alors qu'un ou plusieurs objets du fichier ait un champ : `{
    properties : nombre reacteurs : 4
}`
et que le filtre soit sur le champ `nombre reacteurs` avec la valeur `4`.

Exemple d'un objet en GeoJSON :
```json
{
    "type": "tracks",
    "data": {
        "type": "FeatureCollection",
        "features": [
        {
            "type": "Feature",
            "id": "track.1",
            "geometry": {
            "type": "Point",
            "coordinates": [48.5807528, 1.529444444444443]
            },
            "properties": { //certaines propiétés sont obligatoires
            "altitude": 9600, //obligatoire
            "cap": 145.71,
            "color": { "r": 255, "g": 220, "b": 0 },
            "date": "2021-11-30 16:15:36",  //convertit en timeStampNS
            "identite": "FRIEND",
            "isSelected": false,
            "latitude": 48.5807528, //obligatoire
            "longitude": 1.529444444444443, //obligatoire
            "icao": "1234",
            "type": "a380",
            "headingDeg": 120,
            "pitchDeg": 0,
            "rollDeg": 0,
            "modelRotationDeg": 80
            }
        }
      ]
  }
}
```
Example d'un Filtre en JSON avec les propriété filtrable:
```json
[
  {
    "type": "filters",
    "data": {
      "features": [
        {
          "filtertype": "staticEquals",
          "name": "FRIEND",
          "properties": {
            "field": "identite",
            "query": "FRIEND"
          }
        },
      	{
          "filtertype": "DoubleSlider",
          "name": "Altitude",
          "properties": {
            "field": "altitude",
            "slidmin": 8000,
            "slidmax":10000,
            "step":1 ,
            "value":[7500,10000]
          }
        },
      ]
    }
  },
  {
    "type": "filteredProp",
    "data": {
      "features": [{
        "properties": [
            "altitude",
            "identite",

        ]
        }
      ]
    }
  }
]
```

Les propriété filtrables sont obligatoire afin de pouvoir filtrer de facon générique.
L'utilisateur doit obligatoirement inserer ce fichier (qui vient a etre ajouter directement avec le fichier de filtre si l'utilisateur vient sauvegarder ses filtres).
Ce type de fichier a été penser afin de rendre l'application un maximum générqiue, grâce a ca l'utilisateur peut filtrer sur absolument ce qu'il veut.

Les propriétés obligatoires sont mélangées avec les propiétés génériques afin d'être plus simple d'utilisation.

## Architecture Front
Le front est coder en react et typescript. Une poursuite de langage deja utilisé ici. Il est efficace d'utilisé du typescript car le langage de typage est utile pour comprendre pleins de chose lors de l'apprentissage.
L'architecture est basé sur les normes front/react (UX ?)

## Build Drone

La creation d'un .drone.yml sert a lancé un build drone a chaque commit. Un Dockerfile est egalement creer afin de creer une image docker qui sera recuperer une fois le build drone fini et validé.

.drone.yml
```yml
kind: pipeline
type: docker
name: front-offline Build

workspace:
  base: /js
  path: src//${DRONE_REPO}

steps:
  - name: build_front
    image: docker/base/node:15.3-ssl
    commands:
      - CYPRESS_INSTALL_BINARY=0 npm install
        --registry=http://src/repository/npm/
      - npm run build

  - name: get_tiles
    image: docker/base/node:15.3-ssl
    commands:
      - mkdir tiles && cd tiles
      # -nH --cut-dirs=3 prevents copying tiles' parents folders
      - wget -r -nH --cut-dirs=3
        ftp://ad/repository/cartography/tiles

  - name: publish_front
    image: docker/plugins/docker:22.06-ruche
    settings:
      repo: docker.seed/${DRONE_REPO_OWNER}/seed-offline
      repo_dev: docker.seed/${DRONE_REPO_OWNER}/seed-offline
      context: .
      purge: true
      tags:
        - ${DRONE_COMMIT:0:10}_${DRONE_BRANCH}
        - ${DRONE_TAG}
    environment:
      PLUGIN_USERNAME:
        from_secret: DOCKER_USERNAME
      PLUGIN_PASSWORD:
        from_secret: DOCKER_PASSWORD
    volumes:
      - name: docker-sock
        path: /var/run/docker.sock

volumes:
  - name: docker-sock
    host:
      path: /var/run/docker.sock

```

Fichier Dockerfile
```dockerfile
FROM docker.la-ruche.fr/base/nginx
COPY build /usr/share/nginx/html
COPY tiles /usr/share/nginx/html/tiles
EXPOSE 80
```

## Deploiement

Deploiement du front dans un cluster kubernetes

Explication : Url externe envoie requete (du https) a un INGRESS (porte d'entrée du cluster) qui lui va communiquer a un service (tapant sur le port en question). Le service voit si le pod repond et revoit tout vers l'ingress qui renvoit la reponse. Les communications dans le cluster se font en http, nul besoin d'utiliser des protocoles de securité car communication interne, https est requit des qu'on communique avec l'exterieur.
(faire schéma a l'aide de la photo telephone)

Fichier de deploiement
```yaml
{{- if .Values.frontoffline.enabled }}
# FRONT CONFIG MAP
kind: ConfigMap
apiVersion: v1
metadata:
  name: frontoffline-config
  selfLink: /api/v1/namespaces/default/configmaps/front-config
data:
  config.json: |-
    {
      "imagery":[
        {
          "name": "Natural Earth II",
          "tooltip": "Natural Earth II",
          "tilesUrl": "tiles/NaturalEarthII",
          "layerName": "Shape:DarkEverything",
          "iconUrl": "images/natural_earth_2.png"
        }
      ],
      "terrain": [
        {
          "name": "WGS84 Ellipsoid",
          "tooltip": "WGS84 Ellipsoid",
          "iconUrl": "images/ellipsoid.png",
          "typeTerrain": "terrainEllipsoid"
        }
      ]
    }
---
# FRONT SERVICE
kind: Service
apiVersion: v1
metadata:
  labels:
    app: frontoffline
    version: {{ .Release.Name }}
  name: frontoffline
spec:
  type: ClusterIP
  ports:
  - port: {{ .Values.frontoffline.frontservice.port }}
    targetPort: {{ .Values.frontoffline.frontservice.targetPort }}
    name: frontoffline-ihm
  selector:
    app: frontoffline
---
kind: Ingress
apiVersion: v1
metadata:
  name:  frontoffline-ihm
  annotations:
    nginx.ingress.kubernetes.io/cors-allow-origin: '*'
    nginx.ingress.kubernetes.io/enable-cors: 'true'
spec:
  rules:
  - host: seed-offline.{{ .Values.domain }}
    http:
      paths:
      - pathType: Prefix
        path: /
        backend:
          service:
            name: frontoffline
            port:
              number: 80
---
# FRONT DEPLOYMENT
kind: Deployment
apiVersion: apps/v1
metadata:
  name: frontoffline-deploy
  labels:
    app: frontoffline
    version: {{ .Release.Name }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: frontoffline
  template:
    metadata:
      labels:
        app: frontoffline
    spec:
      {{- with .Values.frontoffline.pullSecret.name }}
      imagePullSecrets:
      - name: {{ . }}
      {{- end }}
      containers:
      - name: frontoffline
        image: {{ .Values.frontoffline.images.frontdeploy }}
        ports:
        - containerPort: 80
        volumeMounts:
        - name: config-json-vol-frontoffline
          mountPath: /usr/share/nginx/html/conf/config.json
          subPath: config.json
      volumes:
      - name: config-json-vol-frontoffline
        configMap:
          name: frontoffline-config
          items:
          - key: config.json
            path: config.json
---
{{- end }}

# VALUES's FILE

# -------- Front Offline -------- #
domain: ihm.com

frontoffline:
  enabled: true
  frontservice:
    port: 80
    targetPort: 80
  pullSecret:
    name: regcred
  images:
    frontdeploy: docker.seed-offline/seed-offline:dockerImage
  appsSettings:
    externalGeoserverAddress: seed.thalesgroup.com/geoserver
    layers:
      DarkEverything: true
      PlanetObserver: true
      OpenStreetMap: true




apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
    containers:
    - name: nginx
      image: nginx:1.14.2
      ports:
      - containerPort: 80

```

## Data

Le principe est de DragNDrop un fichier JSON dans l'application. L'application vient detecter qu'un fichier est sur le point d'etre deposé et affiche une case grise montrant qu'on peut deposer un fichier.
Reagrde si fichier JSON, si Fichier JSON alors le parcours et envoie la data dans un context (explication context)
Utilisation du context car facilité a utilisé la données partout dans l'app.
La data est alors envoyé dans le composant d'affichage qui s'alligne sur le format cesium et permet un affichage carto 3D

Les items sont parsés et reformés dans un model precis pour obtenir toutes les propriétés genériques qui resserviront plus tard (elles sont stocké dans le context)

Le fichier est lui aussi sauvegardé dans un context avec la donnée contenu dedans afin de pouvoir supprimer le fichier a tout moment.

context unit :
```ts
type Context = {
  sensorsData: Map<string, Sensor>;
  setSensorsData: React.Dispatch<React.SetStateAction<Map<string, Sensor>>>;
  plotsData: Map<string, Plot>;
  setPlotsData: React.Dispatch<React.SetStateAction<Map<string, Plot>>>;
  tracksData: Map<string, Track>;
  setTracksData: React.Dispatch<React.SetStateAction<Map<string, Track>>>;
  areasData: Map<string, Area>;
  setAreasData: React.Dispatch<React.SetStateAction<Map<string, Area>>>;
};

type ProviderProps = {
  children: ReactNode;
};

const UnitContext = createContext<Context | undefined>(undefined);

function UnitProvider({ children }: ProviderProps) {
  const [sensorsData, setSensorsData] = useState<Map<string, Sensor>>(
    new Map(),
  );
  const [plotsData, setPlotsData] = useState<Map<string, Plot>>(new Map());
  const [tracksData, setTracksData] = useState<Map<string, Track>>(new Map());
  const [areasData, setAreasData] = useState<Map<string, Area>>(new Map());

  const value: Context = useMemo(
    () => ({
      sensorsData,
      setSensorsData,
      plotsData,
      setPlotsData,
      tracksData,
      setTracksData,
      areasData,
      setAreasData,
    }),
    [sensorsData, plotsData, tracksData, areasData],
  );

  return <UnitContext.Provider value={value}>{children}</UnitContext.Provider>;
}

function useUnits() {
  const context = useContext(UnitContext);
  if (!context) {
    throw new Error('useConfig must be within a UnitContext');
  }

  return context;
}

export { UnitProvider, useUnits };
```

``` ts
const addFilter = (newFilter: Map<string, Filter>) => {
    const { filters, setFilter }: Context = context;
    if (filters && filters.size > 0) {
      setFilter(new Map([...filters, ...newFilter]));
    } else {
      setFilter(newFilter);
    }

    return { filters, setFilter };
  };

  const delFilter = (filtername: string) => {
    const { filters, setFilter }: Context = context;
    filters.delete(filtername);
    setFilter(new Map([...filters]));
    return { filters, setFilter };
  };
```
des states sont rempli avec la data correspondantes
le choix d'avoir separé les differentes items est fait car process different a plusieurs endroit sur les filtres


Mettre fonctionnement du parse fichier, drag n drop, ajout dans mapbeesium (on renomme bien entendu), affichage sur carto + compteur données.
Ne pas hésiter a mettre le context pour montrer fonctionnement



## filtres

Il existe plusieurs type de filtre, staticEquals (on/off), le contains, les double slider et les slider (a l'heure actuelle)
l'implementation de nouveaux filtres est relativement simple car ca ete reflechis dans le code

Un panneau a part entierre est creer car c'est la fonctionnalité premiere de l'application, une séparation et un affichage de chaque type de filtre a lieu pour etre plus clair.
Le fonctionnement est le meme, ils sont non independants, cad que si l'un cache une donnée les autres le prennent en compte et ne viennent pas creer un bug en modifant son etat d'affichage.

Un filtre de couleur est egalement implemanter afin de pouvoir coloriser les items de notre choix. Les filtres fonctionnent exactement comme ceux d'affichage mais agissent sur le parametre de couleur et non d'affichage

photo visu panneau filtre


## Tote

Tote s'active a chaque clic sur un objet pour ressortir toutes les informations de ces objets.
Utile pour savoir les prop de chaque objet, sa couleur d'origine et sa couleur actuelle

montre creation Tote avec informations et rendu

map3D:
```ts
export default function MapBeesium3D() {
  const { config } = useConfig();
  const { sensorsData, plotsData, tracksData, areasData } = useUnits();
  const [display, setDisplay] = useState<boolean>(false);
  const [toteItem, setToteItem] = useState<ToteItem>({
    unit: undefined,
    icon: '',
  });

  if (!config) return null;

  function displayTote(itemClicked: itemClicked) {
    let tote: ToteItem = {
      unit: undefined,
      icon: '',
    };
    if (itemClicked.item !== undefined) {
      setDisplay(true);
      switch (itemClicked.item.type) {
        case Items.sensors:
          return Array.from(sensorsData)
            .filter(([key]) => key === itemClicked.item.id)
            .map(([, value]) => {
              return setToteItem(
                (tote = {
                  unit: value,
                  icon: 'radar', //fix const name in ux/Tote
                }),
              );
            });

        case Items.plots:
          return Array.from(plotsData)
            .filter(([key]) => key === itemClicked.item.id)
            .map(([, value]) => {
              return setToteItem(
                (tote = {
                  unit: value,
                  icon: value.objectType.slice(0, -1),
                }),
              );
            });

        case Items.tracks:
          return Array.from(tracksData)
            .filter(([key]) => key === itemClicked.item.id)
            .map(([, value]) => {
              return setToteItem(
                (tote = {
                  unit: value,
                  icon: value.objectType.slice(0, -1),
                }),
              );
            });
      }
    } else {
      setToteItem(tote);
      setDisplay(false);
    }
  }

  return (
    <div className="offline-map-container">
      <Map3D
        imageryModels={config.imagery}
        terrainModels={config.terrain}
        sensorsCollection={sensorsData}
        tracksCollection={tracksData}
        areasOfInterestCollection={areasData}
        plotsCollection={plotsData}
        hscale={5}
      >
        {/* <Comets points={[...plots.values()]} /> */}

        <UserEventManager
          onLeftClick={(e: any) => {
            displayTote(e);
          }}
          // onMouseMove={(e: any) => {
          //   console.log(e);
          // }}
        />

        <PointPrimitivesCollection3D itemsCollection={plotsData} />
      </Map3D>
      <Draggable defaultPos={{ x: '10%', y: '100px' }}>
        <ToteObjects
          display={display}
          setDisplay={setDisplay}
          toteItem={toteItem}
        />
      </Draggable>
    </div>
  );
}
```
---
# Cette partie technique pour le back (mcd et mld a mettre egalement)

## #TODO

``` graphql
type Sensor {
  uuid : String! @id
  name: String
  trackDetected: [Track] @hasInverse(field:sensorInformation)
}

type Track{
  uuid : String! @id
  name: String
  sensorInformation: [Sensor]
}
```

``` go
type Sensor struct {
	Uuid           string   `json:"Sensor.uuid,omitempty"`
	Name           string   `json:"Sensor.name,omitempty"`
	Track_detected []Track  `json:"Sensor.track_detected,omitempty"`
	DType          []string `json:"dgraph.type,omitempty"`
}
func WriteExampleSensor(ctx context.Context, dgraphClient *dgo.Dgraph) {
	var sensor Sensor
	sensor = Sensor{
		Uuid:  "sensor1",
		Name:  "SensorExample",
		DType: []string{"Sensor"},
	}
	sensorJSON, err := json.Marshal(sensor)
	if err != nil {
		log.Fatal(err)
	}
	mutation := &api.Mutation{}
	mutation.SetJson = sensorJSON
	txn := dgraphClient.NewTxn()
	defer txn.Discard(ctx)
	_, err = txn.Mutate(ctx, mutation)
	if err != nil {
		log.Fatal(err)
	}
	err = txn.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}
}




func init() {
d, err := grpc.Dial("ygt-dgraph-dgraph-alpha:9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    dgraphClient = dgo.NewDgraphClient(api.NewDgraphClient(d))
}
func main() {
    ctx := context.Background()

    client.WriteExampleSensor(ctx, dgraphClient)
    client.WriteExampleTrack(ctx, dgraphClient)
}



type Sensor struct {
  Uuid           string   `json:"Sensor.uuid,omitempty"`
  Name           string   `json:"Sensor.name,omitempty"`
  DType          []string `json:"dgraph.type,omitempty"`
}
func WriteExampleSensor(ctx context.Context, dgraphClient *dgo.Dgraph) {
  var sensor Sensor
  sensor = Sensor{
    Uuid:  "sensor1",
    Name:  "SensorExample",
    DType: []string{"Sensor"},
  }
  sensorJSON, err := json.Marshal(sensor)
  if err != nil {
    log.Fatal(err)
  }
  mutation := &api.Mutation{}
  mutation.SetJson = sensorJSON
  txn := dgraphClient.NewTxn()
  defer txn.Discard(ctx)
  _, err = txn.Mutate(ctx, mutation)
  if err != nil {
    log.Fatal(err)
  }
  err = txn.Commit(ctx)
  if err != nil {
    log.Fatal(err)
  }
}

type Sensor {
    uuid: String! @id @search
    name: String
    latitude_dd: Float
    longitude_dd: Float
    altitude_amsl_m: Float
    track_detected: [Track] @hasInverse(field: sensorsInformation)
    event: [Event] @hasInverse(field: sensorsInformation)
}
type Zone {
    uuid: String! @id @search
    name: String
    sensorsInformation: [Sensor]
    track_associated: [Track]
    event: [Event] @hasInverse(field: zone_associated)
}
type Track {
    uuid: String! @id @search
    name: String
    latitude_dd: Float
    longitude_dd: Float
    altitude_amsl_m: Float
    sensorsInformation: [Sensor]
    event: [Event] @hasInverse(field: track_associated)
}
type Event {
    uuid: String! @id @search
    name: String
    description: String
    latitude_dd: Float
    longitude_dd: Float
    altitude_amsl_m: Float
    functionID: String
    sensorsInformation: [Sensor]
    zone_associated: [Zone]
    track_associated: [Track]
}

func (serv *server) LoadData(addr string) {
    log.Debug("{0} => Waiting for nats datum")

    for buffer := range serv.input {
        switch serv.modelType {
        case "track":
            log.Warn("Decoding track")
            var track models.Units
            err := proto.Unmarshal(buffer, &track)
            if err != nil {
                log.Errorf("We can't unmarshall the track, the error is %v", err)
                continue
            }
            client.WriteTrackGraphQL(addr, track)
        case "event":
            log.Warn("Decoding event")
            var event models.Event
            err := proto.Unmarshal(buffer, &event)
            if err != nil {
                log.Errorf("We can't unmarshall the event, the error is %v", err)
                continue
            }
            client.WriteEvents(addr, event)
        default:
            log.Warnf("Decoding UNKNOWN : %s", serv.modelType)
        }
    }
}


 ```
``` yaml
{{- if .Values.beesplaying.trackarchiver.enabled }}
# BEESPLAYING DEPLOYMENT
kind: Deployment
apiVersion: apps/v1
metadata:
  name: beesplaying-deploy
  labels:
    app: trackarchiver
spec:
  replicas: {{ .Values.beesplaying.trackarchiver.replicas }}
  selector:
    matchLabels:
      app: trackarchiver
  template:
    metadata:
      labels:
        app: trackarchiver
    spec:
      imagePullSecrets:
      - name: {{ .Values.beesplaying.images.pullSecret }}
      {{- with .Values.beesplaying.trackarchiver.nodeSelector }}
      nodeSelector:
{{ toYaml . | indent 8 }}
      {{- end}}
      containers:
      - name: trackarchiver
        image: {{ .Values.beesplaying.images.trackarchiver }}
        imagePullPolicy: IfNotPresent
        args:
          - --nats_address={{ .Values.beesplaying.appsSettings.natsAddress }}
          - --consummer_group={{ .Values.beesplaying.trackarchiver.natsGroup }}
          - --topic_in={{ .Values.beesplaying.trackarchiver.natsTopicIn }}
          - --model_type={{ .Values.beesplaying.trackarchiver.modelType }}
          - --verbose_level={{ .Values.beesplaying.trackarchiver.verboseLevel }}
          - --graphlql_address={{ .Values.beesplaying.appsSettings.graphqlAddress }}
---
{{- end }}
# -------- Track_archiver -------- #
  trackarchiver:
    enabled: true
    replicas: 1
    verboseLevel: 5
    nodeSelector:
      nodetype: sesame
    natsTopicIn: topic_adsb_tracks
    natsGroup: group_adsbeebop
    modelType: track

```
