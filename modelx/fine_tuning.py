import argparse
import torch.cuda
import pandas as pd
import matplotlib.pyplot as plt
from datasets import Dataset
from typing import List, Tuple
from sentence_transformers import (SentenceTransformer, losses, SentenceTransformerTrainingArguments)
from sentence_transformers.trainer import SentenceTransformerTrainer
from sklearn.metrics.pairwise import cosine_similarity
from sklearn.metrics import precision_score, recall_score, f1_score, accuracy_score

device = 'cpu'
if torch.cuda.is_available():
    device = 'cuda'


def calc_cosine_similarity(model: SentenceTransformer, sentence1: str, sentence2: str) -> float:
    # convert the sentences to embeddings
    embeddings = model.encode([sentence1, sentence2])

    # calculate the cosine similarity
    similarities = cosine_similarity([embeddings[0]], [embeddings[1]])

    return similarities[0][0]


def calc_labels(model: SentenceTransformer, test_dataset: Dataset, threshold: float = 0.7) -> Tuple[List[int], List[int]]:
    true_labels = []
    predicted_labels = []

    for row in test_dataset:
        sentence1 = row['texts'][0]
        sentence2 = row['texts'][1]
        true_score = row['label']

        # predict the similarity using the fine-tuned model
        predicted_score = calc_cosine_similarity(model, sentence1, sentence2)

        # convert scores to binary labels based on the threshold
        true_labels.append(1 if true_score >= threshold else 0)
        predicted_labels.append(1 if predicted_score >= threshold else 0)

    return true_labels, predicted_labels


def evaluate(model: SentenceTransformer, test_dataset: Dataset) -> pd.DataFrame:
    true_labels, predicted_labels = calc_labels(model, test_dataset, 0.7)

    accuracy = accuracy_score(true_labels, predicted_labels)
    precision = precision_score(true_labels, predicted_labels)
    recall = recall_score(true_labels, predicted_labels)
    f1 = f1_score(true_labels, predicted_labels)

    # create a DataFrame to display the metrics
    df = pd.DataFrame({
        'Accuracy': [accuracy],
        'Precision': [precision],
        'Recall': [recall],
        'F1_Score': [f1]
    })

    return df


def fine_tuning(csv_file: str, model_path: str, version: str) -> None:
    # load pretrained model.
    model = SentenceTransformer(model_path, device=device)

    # create train dataset.
    rows = Dataset.from_csv(csv_file)
    data = {"guid": [], "texts": [], "label": []}
    for row in rows:
        data["guid"].append(str(row['guid']))
        data["texts"].append([row['sentence1'], row['sentence2']])
        data["label"].append(float(row['score']))

    data = Dataset.from_dict(data)
    dataset = data.train_test_split(test_size=0.2, seed=42)

    train_dataset = dataset['train']
    test_dataset = dataset['test']

    # using cosine similarity as loss function
    train_loss = losses.CosineSimilarityLoss(model=model)

    # set up the training parameters
    training_args = SentenceTransformerTrainingArguments(
        output_dir=".",
        per_device_train_batch_size=16,
        learning_rate=1e-5,  # 0.00001
        num_train_epochs=2,
        weight_decay=0.01,
        warmup_steps=int(len(train_dataset) / 16 * 0.1),
        save_steps=1000,
        logging_dir='./logs',
        logging_steps=100,
    )

    # track the loss at each step
    training_losses = []

    # training loop
    trainer = SentenceTransformerTrainer(
        model=model,
        args=training_args,
        train_dataset=train_dataset,
        loss=train_loss,
    )

    # custom callback to log loss at each step
    def training_step_with_loss_logging(model, inputs):
        loss = SentenceTransformerTrainer.training_step(trainer, model, inputs)
        training_losses.append(loss.item())
        return loss

    trainer.training_step = training_step_with_loss_logging

    # start train
    trainer.train()

    # Plot the training loss after training
    plt.plot(training_losses)
    plt.title('Training Loss Over Time')
    plt.xlabel('Training Steps')
    plt.ylabel('Loss')
    plt.savefig("training_loss.png")

    # save the fine-tuned model
    ft_model_path = model_path + '_FT_' + version
    model.save(ft_model_path)

    # evaluate on the test set
    origin_model = SentenceTransformer(model_path, device=device)
    df = evaluate(origin_model, test_dataset)
    print("origin model evaluation")
    print(df)

    print("")

    df = evaluate(model, test_dataset)
    print("fine-tuned model evaluation")
    print(df)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Modelx fine tuning")
    parser.add_argument('--dataset', type=str, required=True, help="Path to the dataset CSV file.")
    parser.add_argument('--model', type=str, required=True, help="Path to the model folder.")
    parser.add_argument('--version', type=str, required=True, help="Set the version of fine tuning model.")
    args = parser.parse_args()

    if args.dataset and args.model and args.version:
        fine_tuning(args.dataset, args.model, args.version)
